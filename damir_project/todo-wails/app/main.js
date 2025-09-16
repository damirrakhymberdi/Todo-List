// Импорт Wails runtime
import { Get } from '@wails/runtime';

// Получение бэкенда
let backend;

// Инициализация при загрузке
document.addEventListener('DOMContentLoaded', async () => {
    try {
        // Получаем бэкенд из Wails
        backend = await Get('backend.App');
        await refreshTasks();
    } catch (error) {
        console.error('Ошибка инициализации:', error);
        showError('Ошибка инициализации приложения');
    }
});

// Элементы DOM
const elements = {
    form: document.getElementById('addForm'),
    title: document.getElementById('title'),
    due: document.getElementById('due'),
    priority: document.getElementById('priority'),
    filter: document.getElementById('filter'),
    sort: document.getElementById('sort'),
    container: document.getElementById('tasksContainer'),
    themeBtn: document.getElementById('themeBtn')
};

// Обработчики событий
elements.form.addEventListener('submit', handleAddTask);
elements.filter.addEventListener('change', refreshTasks);
elements.sort.addEventListener('change', refreshTasks);
elements.themeBtn.addEventListener('click', toggleTheme);

// Обработчик добавления задачи
async function handleAddTask(e) {
    e.preventDefault();

    const title = elements.title.value.trim();
    if (!title) {
        showError('Введите название задачи');
        return;
    }

    try {
        const dueISO = elements.due.value ? new Date(elements.due.value).toISOString() : null;
        await backend.UC.Add(title, dueISO, elements.priority.value);

        // Очистка формы
        elements.title.value = '';
        elements.due.value = '';
        elements.priority.value = 'medium';

        await refreshTasks();
    } catch (error) {
        console.error('Ошибка добавления задачи:', error);
        showError('Ошибка добавления задачи');
    }
}

// Обработчик переключения статуса задачи
async function toggleTask(id) {
    try {
        await backend.UC.ToggleDone(id);
        await refreshTasks();
    } catch (error) {
        console.error('Ошибка переключения статуса:', error);
        showError('Ошибка обновления задачи');
    }
}

// Обработчик удаления задачи
async function deleteTask(id) {
    if (confirm('Вы уверены, что хотите удалить эту задачу?')) {
        try {
            await backend.UC.Delete(id);
            await refreshTasks();
        } catch (error) {
            console.error('Ошибка удаления задачи:', error);
            showError('Ошибка удаления задачи');
        }
    }
}

// Обновление списка задач
async function refreshTasks() {
    try {
        elements.container.innerHTML = '<div class="loading">Загрузка...</div>';

        const tasks = await backend.UC.List(elements.filter.value, elements.sort.value);

        if (tasks.length === 0) {
            elements.container.innerHTML = '<div class="loading">Нет задач</div>';
            return;
        }

        // Разделение на активные и выполненные
        const activeTasks = tasks.filter(task => !task.done);
        const doneTasks = tasks.filter(task => task.done);

        let html = '';

        // Активные задачи
        if (activeTasks.length > 0 && elements.filter.value !== 'done') {
            html += '<div class="task-section">';
            html += '<h3>🟢 Активные задачи</h3>';
            html += '<ul class="task-list">';
            activeTasks.forEach(task => {
                html += createTaskHTML(task);
            });
            html += '</ul></div>';
        }

        // Выполненные задачи
        if (doneTasks.length > 0 && elements.filter.value !== 'active') {
            html += '<div class="task-section">';
            html += '<h3>✅ Выполненные задачи</h3>';
            html += '<ul class="task-list">';
            doneTasks.forEach(task => {
                html += createTaskHTML(task);
            });
            html += '</ul></div>';
        }

        elements.container.innerHTML = html;

        // Добавление обработчиков событий
        addEventListeners();

    } catch (error) {
        console.error('Ошибка загрузки задач:', error);
        showError('Ошибка загрузки задач');
    }
}

// Создание HTML для задачи
function createTaskHTML(task) {
    const dueDate = task.dueAt ? new Date(task.dueAt).toLocaleString('ru-RU') : '';
    const priorityClass = task.priority.toLowerCase();
    const priorityText = {
        'high': '🔴 Высокий',
        'medium': '🟡 Средний',
        'low': '🟢 Низкий'
    }[task.priority] || task.priority;

    return `
        <li class="task-item ${task.done ? 'done' : ''} priority-${priorityClass}">
            <input type="checkbox" 
                   class="task-checkbox" 
                   ${task.done ? 'checked' : ''} 
                   onchange="toggleTask('${task.id}')">
            <div class="task-content">
                <div class="task-title ${task.done ? 'done' : ''}">${escapeHtml(task.title)}</div>
                <div class="task-meta">
                    <span class="task-priority ${priorityClass}">${priorityText}</span>
                    ${dueDate ? `<span class="task-due">📅 ${dueDate}</span>` : ''}
                </div>
            </div>
            <div class="task-actions">
                <button class="delete-btn" onclick="deleteTask('${task.id}')">🗑️ Удалить</button>
            </div>
        </li>
    `;
}

// Добавление обработчиков событий
function addEventListeners() {
    // Обработчики уже добавлены через onclick в HTML
}

// Переключение темы
function toggleTheme() {
    document.body.classList.toggle('dark');
    const isDark = document.body.classList.contains('dark');
    elements.themeBtn.textContent = isDark ? '☀️' : '🌙';
    localStorage.setItem('theme', isDark ? 'dark' : 'light');
}

// Загрузка сохраненной темы
function loadTheme() {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'dark') {
        document.body.classList.add('dark');
        elements.themeBtn.textContent = '☀️';
    }
}

// Показ ошибки
function showError(message) {
    elements.container.innerHTML = `<div class="loading" style="color: #e74c3c;">❌ ${message}</div>`;
}

// Экранирование HTML
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Инициализация темы
loadTheme();
