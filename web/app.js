const state = {
  notes: [],
  activeId: null,
  query: "",
};

const els = {
  noteList: document.getElementById("note-list"),
  listEmpty: document.getElementById("list-empty"),
  searchInput: document.getElementById("search-input"),
  noteForm: document.getElementById("note-form"),
  titleInput: document.getElementById("title-input"),
  contentInput: document.getElementById("content-input"),
  editorTitle: document.getElementById("editor-title"),
  saveBtn: document.getElementById("save-btn"),
  deleteBtn: document.getElementById("delete-btn"),
  cancelBtn: document.getElementById("cancel-btn"),
  newNoteBtn: document.getElementById("new-note-btn"),
  statusMsg: document.getElementById("status-msg"),
};

const timeFormatter = new Intl.DateTimeFormat("zh-CN", {
  hour12: false,
  year: "numeric",
  month: "2-digit",
  day: "2-digit",
  hour: "2-digit",
  minute: "2-digit",
});

boot();

function boot() {
  bindEvents();
  void refreshNotes();
}

function bindEvents() {
  els.noteForm.addEventListener("submit", (event) => {
    event.preventDefault();
    void saveNote();
  });

  els.newNoteBtn.addEventListener("click", () => {
    switchToCreateMode();
    showStatus("已进入新建模式", "ok");
  });

  els.cancelBtn.addEventListener("click", () => {
    if (state.activeId !== null) {
      const note = state.notes.find((item) => item.id === state.activeId);
      if (note) {
        selectNote(note.id);
        showStatus("已取消修改", "ok");
        return;
      }
    }
    switchToCreateMode();
    showStatus("已清空编辑器", "ok");
  });

  els.deleteBtn.addEventListener("click", () => {
    void removeNote();
  });

  els.searchInput.addEventListener("input", (event) => {
    state.query = event.target.value.trim().toLowerCase();
    renderNoteList();
  });
}

async function refreshNotes() {
  try {
    const notes = await requestJSON("/notes");
    state.notes = Array.isArray(notes) ? notes : [];

    renderNoteList();

    if (state.activeId !== null) {
      const active = state.notes.find((item) => item.id === state.activeId);
      if (active) {
        fillEditor(active);
      } else {
        switchToCreateMode();
      }
    }
  } catch (error) {
    showStatus(getErrorMessage(error), "error");
  }
}

function renderNoteList() {
  const keyword = state.query;
  const filtered = state.notes.filter((item) => {
    if (!keyword) return true;
    return (
      item.title.toLowerCase().includes(keyword) ||
      item.content.toLowerCase().includes(keyword)
    );
  });

  els.noteList.innerHTML = "";
  els.listEmpty.classList.toggle("hidden", filtered.length > 0);

  for (const item of filtered) {
    const li = document.createElement("li");
    li.className = "note-item";
    if (item.id === state.activeId) {
      li.classList.add("active");
    }

    const title = document.createElement("p");
    title.className = "note-item-title";
    title.textContent = item.title;

    const meta = document.createElement("p");
    meta.className = "note-item-meta";
    meta.textContent = `更新于 ${formatTime(item.updated_at)}`;

    const body = document.createElement("p");
    body.className = "note-item-body";
    body.textContent = item.content;

    li.appendChild(title);
    li.appendChild(meta);
    li.appendChild(body);

    li.addEventListener("click", () => {
      selectNote(item.id);
    });

    els.noteList.appendChild(li);
  }
}

function selectNote(id) {
  const note = state.notes.find((item) => item.id === id);
  if (!note) return;

  state.activeId = note.id;
  fillEditor(note);
  renderNoteList();
}

function fillEditor(note) {
  els.titleInput.value = note.title;
  els.contentInput.value = note.content;
  els.editorTitle.textContent = `编辑笔记 #${note.id}`;
  els.saveBtn.textContent = "更新笔记";
  els.deleteBtn.classList.remove("hidden");
  els.cancelBtn.classList.remove("hidden");
}

function switchToCreateMode() {
  state.activeId = null;
  els.noteForm.reset();
  els.editorTitle.textContent = "新建笔记";
  els.saveBtn.textContent = "保存笔记";
  els.deleteBtn.classList.add("hidden");
  els.cancelBtn.classList.add("hidden");
  renderNoteList();
}

async function saveNote() {
  const payload = {
    title: els.titleInput.value.trim(),
    content: els.contentInput.value.trim(),
  };

  if (!payload.title) {
    showStatus("标题不能为空", "error");
    return;
  }
  if (!payload.content) {
    showStatus("内容不能为空", "error");
    return;
  }

  setSaving(true);
  try {
    let saved;
    if (state.activeId === null) {
      saved = await requestJSON("/notes", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      showStatus("创建成功", "ok");
    } else {
      saved = await requestJSON(`/notes/${state.activeId}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      showStatus("更新成功", "ok");
    }

    state.activeId = saved.id;
    await refreshNotes();
  } catch (error) {
    showStatus(getErrorMessage(error), "error");
  } finally {
    setSaving(false);
  }
}

async function removeNote() {
  if (state.activeId === null) return;

  const sure = window.confirm("确认删除这条笔记吗？");
  if (!sure) return;

  setSaving(true);
  try {
    await requestJSON(`/notes/${state.activeId}`, { method: "DELETE" }, true);
    switchToCreateMode();
    await refreshNotes();
    showStatus("删除成功", "ok");
  } catch (error) {
    showStatus(getErrorMessage(error), "error");
  } finally {
    setSaving(false);
  }
}

function setSaving(isSaving) {
  els.saveBtn.disabled = isSaving;
  els.deleteBtn.disabled = isSaving;
  els.newNoteBtn.disabled = isSaving;
}

async function requestJSON(url, options = {}, allowNoContent = false) {
  const response = await fetch(url, options);
  if (!response.ok) {
    let message = "请求失败";
    try {
      const data = await response.json();
      if (data && typeof data.error === "string") {
        message = data.error;
      }
    } catch (_) {
      // ignore parse error
    }
    throw new Error(message);
  }

  if (allowNoContent && response.status === 204) {
    return null;
  }
  return response.json();
}

function showStatus(text, type = "") {
  els.statusMsg.textContent = text;
  els.statusMsg.classList.remove("ok", "error");
  if (type) {
    els.statusMsg.classList.add(type);
  }
}

function getErrorMessage(error) {
  if (error instanceof Error && error.message) {
    return error.message;
  }
  return "发生未知错误";
}

function formatTime(value) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return timeFormatter.format(date);
}
