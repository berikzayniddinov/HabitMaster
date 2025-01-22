// Создание карточки элемента
function createCard(id, title, description, date, type, extra = '', actions = []) {
    const container = document.getElementById(`${type}-container`);
    const card = document.createElement('div');
    card.classList.add('card');
    card.dataset.name = title || 'No Name'; // Привязка к текущему названию

    card.innerHTML = `
        <h3>${title || 'No Title'}</h3>
        <p>${description || 'No Description'}</p>
        <p><strong>Date:</strong> ${date || 'No Date'}</p>
        ${extra ? `<p>${extra}</p>` : ''}
        <div class="card-actions"></div>
    `;

    const actionsContainer = card.querySelector('.card-actions');
    actions.forEach(action => actionsContainer.appendChild(action));
    container.appendChild(card);
}


// Создание кнопки действия
function createActionButton(text, className, onclick) {
    const button = document.createElement('button');
    button.textContent = text;
    button.classList.add(className);
    button.onclick = onclick;
    return button;
}

// Добавление привычки
let currentHabitPage = 1;

// Получение привычек
async function getHabits() {
    const filter = document.getElementById('habit-filter')?.value || '';
    const sort = document.getElementById('habit-sort')?.value || '';
    const page = currentHabitPage;

    const url = `http://localhost:8080/api/habits?filter=${encodeURIComponent(filter)}&sort=${encodeURIComponent(sort)}&page=${page}`;
    try {
        const response = await fetch(url, { method: 'GET' });
        if (!response.ok) {
            throw new Error(`Error fetching habits: ${response.statusText}`);
        }

        const habits = await response.json();
        const container = document.getElementById('habits-container');
        container.innerHTML = ''; // Очистка контейнера

        if (!habits || habits.length === 0) {
            container.innerHTML = '<p>No habits found.</p>';
            return;
        }

        habits.forEach(habit => {
            const card = document.createElement('div');
            card.className = 'card';
            card.innerHTML = `
                <h3>${habit.name}</h3>
                <p><strong>Description:</strong> ${habit.description}</p>
                <p><strong>Created At:</strong> ${new Date(habit.created_at).toLocaleString()}</p>
                <div class="card-actions">
                    <button class="edit" onclick="editHabit('${habit.name}', '${habit.description}')">Edit</button>
                    <button class="delete" onclick="deleteHabit('${habit.name}')">Delete</button>
                </div>
            `;
            container.appendChild(card);
        });

        updateHabitPaginationControls();
    } catch (error) {
        console.error('Error fetching habits:', error.message);
        alert('Failed to fetch habits.');
    }
}

function applyHabitFilterSort() {
    getHabits();
}

// Пагинация
function updateHabitPaginationControls() {
    const paginationContainer = document.getElementById('habits-pagination');
    paginationContainer.innerHTML = `
        <button onclick="prevHabitPage()">Previous</button>
        <button onclick="nextHabitPage()">Next</button>
    `;
}

function prevHabitPage() {
    if (currentHabitPage > 1) {
        currentHabitPage--;
        getHabits();
    }
}

function nextHabitPage() {
    currentHabitPage++;
    getHabits();
}

// Добавление привычки
async function addHabit() {
    const name = prompt('Enter habit name:');
    const description = prompt('Enter habit description:');

    if (!name) {
        alert('Habit name is required!');
        return;
    }

    const habit = { name, description };

    try {
        const response = await fetch('http://localhost:8080/api/habits', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(habit),
        });

        if (response.ok) {
            alert('Habit successfully added!');
            getHabits();
        } else {
            const errorText = await response.text();
            console.error('Error adding habit:', errorText);
            alert('Error adding habit.');
        }
    } catch (error) {
        console.error('Unexpected error adding habit:', error);
        alert('Unexpected error adding habit.');
    }
}

// Удаление привычки
async function deleteHabit(name) {
    try {
        const response = await fetch(`http://localhost:8080/api/habits?name=${encodeURIComponent(name)}`, {
            method: 'DELETE',
        });

        if (response.ok) {
            alert('Habit successfully deleted!');
            getHabits(); // Обновить список привычек
        } else {
            const errorText = await response.text();
            console.error('Error deleting habit:', errorText);
            alert(`Error deleting habit: ${errorText}`);
        }
    } catch (error) {
        console.error('Unexpected error deleting habit:', error);
        alert('Unexpected error deleting habit.');
    }
}

// Редактирование привычки
async function editHabit(oldName, currentDescription) {
    const newName = prompt('Enter new habit name:', oldName);
    const newDescription = prompt('Enter new habit description:', currentDescription);

    if (!newName) {
        alert('Habit name is required!');
        return;
    }

    const updatedHabit = { oldName, name: newName, description: newDescription };

    try {
        const response = await fetch('http://localhost:8080/api/habits', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(updatedHabit),
        });

        if (response.ok) {
            alert('Habit successfully updated!');
            getHabits();
        } else {
            const errorText = await response.text();
            console.error('Error updating habit:', errorText);
            alert('Error updating habit.');
        }
    } catch (error) {
        console.error('Unexpected error updating habit:', error);
        alert('Unexpected error updating habit.');
    }
}

// Инициализация
window.onload = getHabits;




// Вызов получения привычек при загрузке страницы
window.onload = getHabits;


let currentGoalPage = 1;

// Получение целей
async function getGoals() {
    const filter = document.getElementById('goal-filter')?.value || '';
    const sort = document.getElementById('goal-sort')?.value || '';
    const page = currentGoalPage;

    const url = `http://localhost:8080/api/goals?filter=${encodeURIComponent(filter)}&sort=${encodeURIComponent(sort)}&page=${page}`;
    try {
        const response = await fetch(url, { method: 'GET' });
        if (!response.ok) {
            throw new Error(`Error fetching goals: ${response.statusText}`);
        }

        const goals = await response.json();
        const container = document.getElementById('goals-container');
        container.innerHTML = ''; // Очистка контейнера

        if (!goals || goals.length === 0) {
            container.innerHTML = '<p>No goals found.</p>';
            return;
        }

        goals.forEach(goal => {
            const card = document.createElement('div');
            card.className = 'card';
            card.innerHTML = `
                <h3>${goal.name}</h3>
                <p><strong>Description:</strong> ${goal.description}</p>
                <p><strong>Deadline:</strong> ${new Date(goal.deadline).toLocaleDateString()}</p>
                <div class="card-actions">
                    <button class="edit" onclick="editGoal('${goal.name}', '${goal.description}', '${goal.deadline}')">Edit</button>
                    <button class="delete" onclick="deleteGoalByName('${goal.name}')">Delete</button>
                </div>
            `;
            container.appendChild(card);
        });

        updateGoalPaginationControls();
    } catch (error) {
        console.error('Error fetching goals:', error.message);
        alert('Failed to fetch goals.');
    }
}

// Применение фильтров и сортировки
function applyGoalFilterSort() {
    getGoals();
}

// Пагинация
function updateGoalPaginationControls() {
    const paginationContainer = document.getElementById('goals-pagination');
    paginationContainer.innerHTML = `
        <button onclick="prevGoalPage()">Previous</button>
        <button onclick="nextGoalPage()">Next</button>
    `;
}

function prevGoalPage() {
    if (currentGoalPage > 1) {
        currentGoalPage--;
        getGoals();
    }
}

function nextGoalPage() {
    currentGoalPage++;
    getGoals();
}

// Добавление цели
async function addGoal() {
    const name = prompt('Enter goal name:');
    const description = prompt('Enter goal description:');
    const deadline = prompt('Enter deadline (YYYY-MM-DD):');

    if (!name || !deadline) {
        alert('Name and deadline are required!');
        return;
    }

    const goal = { name, description, deadline };

    try {
        const response = await fetch('http://localhost:8080/api/goals', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(goal),
        });

        if (response.ok) {
            alert('Goal successfully added!');
            getGoals();
        } else {
            const errorText = await response.text();
            console.error('Error adding goal:', errorText);
            alert('Error adding goal.');
        }
    } catch (error) {
        console.error('Unexpected error adding goal:', error);
        alert('Unexpected error adding goal.');
    }
}

// Удаление цели по имени
async function deleteGoalByName(name) {
    try {
        const response = await fetch(`http://localhost:8080/api/goals/deleteByName?name=${encodeURIComponent(name)}`, {
            method: 'DELETE',
        });

        if (response.ok) {
            alert('Goal successfully deleted!');
            getGoals(); // Обновить список целей
        } else {
            const errorText = await response.text();
            console.error('Error deleting goal:', errorText);
            alert(`Error deleting goal: ${errorText}`);
        }
    } catch (error) {
        console.error('Unexpected error deleting goal:', error);
        alert('Unexpected error deleting goal.');
    }
}




// Редактирование цели
async function editGoal(currentName, currentDescription, currentDeadline) {
    const newName = prompt('Enter new name:', currentName);
    const newDescription = prompt('Enter new description:', currentDescription);
    const newDeadline = prompt('Enter new deadline (YYYY-MM-DD):', currentDeadline);

    if (!newName || !newDeadline) {
        alert('Name and deadline are required!');
        return;
    }

    const updatedGoal = { oldName: currentName, name: newName, description: newDescription, deadline: newDeadline };

    try {
        const response = await fetch('http://localhost:8080/api/goals', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(updatedGoal),
        });

        if (response.ok) {
            alert('Goal successfully updated!');
            getGoals();
        } else {
            const errorText = await response.text();
            console.error('Error updating goal:', errorText);
            alert('Error updating goal.');
        }
    } catch (error) {
        console.error('Unexpected error updating goal:', error);
        alert('Unexpected error updating goal.');
    }
}

// Инициализация
window.onload = getGoals;

// Вызов получения целей при загрузке страницы
window.onload = getGoals;


// Automatically load goals when the page loads
window.onload = getGoals;

window.onload = getGoals();

// Добавление достижения
// Работа с достижениями (Achievements)

// Добавление достижения
// Получение достижений
let currentAchievementPage = 1;

// Получение достижений
async function getAchievements() {
    const filter = document.getElementById('achievement-filter')?.value || '';
    const sort = document.getElementById('achievement-sort')?.value || '';
    const page = currentAchievementPage;

    const url = `http://localhost:8080/api/achievements?filter=${encodeURIComponent(filter)}&sort=${encodeURIComponent(sort)}&page=${page}`;
    try {
        const response = await fetch(url, { method: 'GET' });
        if (!response.ok) {
            throw new Error(`Error fetching achievements: ${response.statusText}`);
        }

        const achievements = await response.json();
        const container = document.getElementById('achievements-container');
        container.innerHTML = ''; // Очистка контейнера

        if (!achievements || achievements.length === 0) {
            container.innerHTML = '<p>No achievements found.</p>';
            return;
        }

        achievements.forEach(achievement => {
            const card = document.createElement('div');
            card.className = 'card';
            card.innerHTML = `
                <h3>${achievement.title}</h3>
                <p><strong>Description:</strong> ${achievement.description}</p>
                <p><strong>Date:</strong> ${new Date(achievement.date).toLocaleString()}</p>
                <div class="card-actions">
                    <button class="edit" onclick="editAchievement('${achievement.title}', '${achievement.description}')">Edit</button>
                    <button class="delete" onclick="deleteAchievement('${achievement.title}')">Delete</button>
                </div>
            `;
            container.appendChild(card);
        });

        updateAchievementPaginationControls();
    } catch (error) {
        console.error('Error fetching achievements:', error.message);
        alert('Failed to fetch achievements.');
    }
}

// Фильтр и сортировка
function applyAchievementFilterSort() {
    getAchievements();
}

// Пагинация
function updateAchievementPaginationControls() {
    const paginationContainer = document.getElementById('achievements-pagination');
    paginationContainer.innerHTML = `
        <button onclick="prevAchievementPage()">Previous</button>
        <button onclick="nextAchievementPage()">Next</button>
    `;
}

function prevAchievementPage() {
    if (currentAchievementPage > 1) {
        currentAchievementPage--;
        getAchievements();
    }
}

function nextAchievementPage() {
    currentAchievementPage++;
    getAchievements();
}

// Добавление достижения
async function addAchievement() {
    const title = prompt('Enter achievement title:');
    const description = prompt('Enter achievement description:');

    if (!title || !description) {
        alert('Both title and description are required!');
        return;
    }

    const achievement = { title, description };

    try {
        const response = await fetch('http://localhost:8080/api/achievements', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(achievement),
        });

        if (response.ok) {
            alert('Achievement successfully added!');
            getAchievements();
        } else {
            const errorText = await response.text();
            console.error('Error adding achievement:', errorText);
            alert('Error adding achievement.');
        }
    } catch (error) {
        console.error('Unexpected error adding achievement:', error);
        alert('Unexpected error adding achievement.');
    }
}

// Удаление достижения
async function deleteAchievement(title) {
    try {
        const response = await fetch(`http://localhost:8080/api/achievements?title=${encodeURIComponent(title)}`, {
            method: 'DELETE',
        });

        if (response.ok) {
            alert('Achievement successfully deleted!');
            getAchievements();
        } else {
            const errorText = await response.text();
            console.error('Error deleting achievement:', errorText);
            alert(`Error deleting achievement: ${errorText}`);
        }
    } catch (error) {
        console.error('Unexpected error deleting achievement:', error);
        alert('Unexpected error deleting achievement.');
    }
}

// Редактирование достижения
async function editAchievement(currentTitle, currentDescription) {
    const newTitle = prompt('Enter new title:', currentTitle);
    const newDescription = prompt('Enter new description:', currentDescription);

    if (!newTitle || !newDescription) {
        alert('Both title and description are required!');
        return;
    }

    const updatedAchievement = { oldTitle: currentTitle, title: newTitle, description: newDescription };

    try {
        const response = await fetch('http://localhost:8080/api/achievements', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(updatedAchievement),
        });

        if (response.ok) {
            alert('Achievement successfully updated!');
            getAchievements();
        } else {
            const errorText = await response.text();
            console.error('Error updating achievement:', errorText);
            alert('Failed to update achievement.');
        }
    } catch (error) {
        console.error('Unexpected error updating achievement:', error);
        alert('Unexpected error updating achievement.');
    }
}

// Инициализация
window.onload = function () {
    getAchievements();
};

let currentNotificationPage = 1;

// Получение уведомлений
async function getNotifications() {
    const filter = document.getElementById('notification-filter')?.value || '';
    const sort = document.getElementById('notification-sort')?.value || '';
    const page = currentNotificationPage;

    const url = `http://localhost:8080/api/notifications?filter=${encodeURIComponent(filter)}&sort=${encodeURIComponent(sort)}&page=${page}`;
    try {
        const response = await fetch(url, { method: 'GET' });
        if (!response.ok) {
            throw new Error(`Error fetching notifications: ${response.statusText}`);
        }

        const notifications = await response.json();
        const container = document.getElementById('notifications-container');
        container.innerHTML = ''; // Очистка контейнера

        if (!notifications || notifications.length === 0) {
            container.innerHTML = '<p>No notifications found.</p>';
            return;
        }

        notifications.forEach(notification => {
            const card = document.createElement('div');
            card.className = 'card';
            card.innerHTML = `
                <h3>${notification.message}</h3>
                <p><strong>Scheduled At:</strong> ${new Date(notification.scheduled_at).toLocaleString()}</p>
                <p><strong>Is Sent:</strong> ${notification.is_sent ? 'Yes' : 'No'}</p>
                <div class="card-actions">
                    <button class="edit" onclick="editNotification(${notification.id}, '${notification.message}', '${notification.scheduled_at}', ${notification.is_sent})">Edit</button>
                    <button class="delete" onclick="deleteNotificationByName('${notification.message}')">Delete</button>
                </div>
            `;
            container.appendChild(card);
        });

        updatePaginationControls();
    } catch (error) {
        console.error('Error fetching notifications:', error.message);
        alert('Failed to fetch notifications.');
    }
}


function applyNotificationFilterSort() {
    getNotifications();
}


// Пагинация
function updatePaginationControls() {
    const paginationContainer = document.getElementById('notifications-pagination');
    paginationContainer.innerHTML = `
        <button onclick="prevPage()">Previous</button>
        <button onclick="nextPage()">Next</button>
    `;
}

function prevPage() {
    if (currentNotificationPage > 1) {
        currentNotificationPage--;
        getNotifications();
    }
}

function nextPage() {
    currentNotificationPage++;
    getNotifications();
}

// Добавление уведомления
async function addNotification() {
    const message = prompt('Enter notification message:');
    const scheduledAt = prompt('Enter scheduled date and time (YYYY-MM-DDTHH:mm:ss):');

    if (!message || !scheduledAt) {
        alert('Both message and scheduled date are required!');
        return;
    }

    const notification = { message, scheduled_at: scheduledAt };

    try {
        const response = await fetch('http://localhost:8080/api/notifications', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(notification),
        });

        if (response.ok) {
            alert('Notification successfully added!');
            getNotifications();
        } else {
            const errorText = await response.text();
            console.error('Error adding notification:', errorText);
            alert('Error adding notification.');
        }
    } catch (error) {
        console.error('Unexpected error adding notification:', error);
        alert('Unexpected error adding notification.');
    }
}

// Удаление уведомления
// Удаление уведомления по имени
async function deleteNotificationByName(message) {
    try {
        const response = await fetch(`http://localhost:8080/api/notifications/deleteByName?message=${encodeURIComponent(message)}`, {
            method: 'DELETE',
        });

        if (response.ok) {
            alert('Notification successfully deleted!');
            getNotifications(); // Обновить список уведомлений
        } else {
            const errorText = await response.text();
            console.error('Error deleting notification:', errorText);
            alert(`Error deleting notification: ${errorText}`);
        }
    } catch (error) {
        console.error('Unexpected error deleting notification:', error);
        alert('Unexpected error deleting notification.');
    }
}



// Редактирование уведомления
async function editNotification(id, currentMessage, scheduledAt, isSent) {
    const newMessage = prompt('Enter new message:', currentMessage);
    const newScheduledAt = prompt('Enter new scheduled date (YYYY-MM-DDTHH:mm:ss):', scheduledAt);

    if (!newMessage || !newScheduledAt) {
        alert('All fields are required!');
        return;
    }

    const updatedNotification = {
        oldMessage: currentMessage,
        message: newMessage,
        scheduled_at: newScheduledAt,
        is_sent: isSent
    };

    try {
        const response = await fetch('http://localhost:8080/api/notifications', {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(updatedNotification),
        });

        if (response.ok) {
            alert('Notification successfully updated!');
            getNotifications();
        } else {
            const errorText = await response.text();
            console.error('Error updating notification:', errorText);
            alert('Failed to update notification.');
        }
    } catch (error) {
        console.error('Unexpected error updating notification:', error);
        alert('Unexpected error updating notification.');
    }
}

// Инициализация
window.onload = getNotifications;







// Вызов получения уведомлений при загрузке страницы
window.onload = getNotifications;


window.onload = function () {
    getHabits();
    getGoals();
    getAchievements();
    getNotifications();
};

async function sendEmail() {
    const recipients = document.getElementById('email-recipients').value;
    const subject = document.getElementById('email-subject').value;
    const body = document.getElementById('email-body').value;
    const attachment = document.getElementById('email-attachment').files[0];

    if (!recipients || !subject || !body) {
        alert('All fields are required!');
        return;
    }

    const formData = new FormData();
    formData.append('recipients', recipients);
    formData.append('subject', subject);
    formData.append('body', body);
    if (attachment) {
        formData.append('attachment', attachment);
    }
    console.log([...formData.entries()]);


    try {
        const response = await fetch('http://localhost:8080/api/admin/send-mass-email', {
            method: 'POST',
            body: formData, // Передача данных через FormData
        });

        if (response.ok) {
            alert('Email sent successfully!');
            document.getElementById('email-recipients').value = '';
            document.getElementById('email-subject').value = '';
            document.getElementById('email-body').value = '';
            document.getElementById('email-attachment').value = '';
        } else {
            const errorText = await response.text();
            console.error('Error sending email:', errorText);
            alert('Failed to send email.');
        }
    } catch (error) {
        console.error('Unexpected error:', error);
        alert('Unexpected error occurred.');
    }
}
document.addEventListener('DOMContentLoaded', () => {
    const registrationForm = document.getElementById('registration-form');
    const getUsersButton = document.getElementById('get-users');
    const userList = document.getElementById('user-list');

    // Функция отображения секции профиля
    function showProfileSection(user) {
        const profileSection = document.getElementById('profile-section');
        const profileNameDisplay = document.getElementById('profile-name-display');
        const profileImage = document.getElementById('profile-image');

        if (profileSection) {
            profileSection.style.display = 'block';
            if (profileNameDisplay) {
                profileNameDisplay.textContent = user.name || "Anonymous"; // Отображение имени
            }
            if (profileImage) {
                profileImage.src = user.profile_picture || 'default-profile.png'; // Установка фото профиля
            }
        }
    }

    // Обработчик отправки формы регистрации
    if (registrationForm) {
        registrationForm.addEventListener('submit', async (event) => {
            event.preventDefault(); // Предотвращаем перезагрузку страницы

            // Получаем данные из формы
            const formData = new FormData(registrationForm);
            const userData = {
                name: formData.get('name'),
                email: formData.get('email'),
                password: formData.get('password'),
            };

            console.log('Registration data:', userData); // Для проверки

            try {
                // Отправляем данные на сервер
                const response = await fetch('/signup', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify(userData),
                });

                console.log('Response:', response); // Проверяем ответ сервера

                // Обрабатываем ответ сервера
                if (response.ok) {
                    const result = await response.json();
                    console.log('Server result:', result); // Для проверки
                    localStorage.setItem("authToken", result.token); // Сохраняем токен
                    alert(result.message); // Сообщение об успехе
                    registrationForm.reset(); // Очищаем форму
                    window.location.href = 'main.html'; // Переход на главную страницу
                } else {
                    const error = await response.text();
                    alert(`Ошибка: ${error}`);
                }
            } catch (err) {
                console.error('Ошибка при регистрации:', err);
                alert('Произошла ошибка при регистрации');
            }
        });
    }

    // Обработчик для кнопки "Получить пользователей"
    if (getUsersButton) {
        getUsersButton.addEventListener('click', async () => {
            try {
                // Запрашиваем список пользователей
                const response = await fetch('/api/get-users');
                if (response.ok) {
                    const users = await response.json();

                    // Очищаем список
                    if (userList) {
                        userList.innerHTML = '';
                    }

                    // Добавляем пользователей в список
                    users.forEach((user) => {
                        const listItem = document.createElement('li');
                        listItem.textContent = `ID: ${user.user_id}, Имя: ${user.name}, Email: ${user.email}`;
                        userList.appendChild(listItem);
                    });
                } else {
                    const error = await response.text();
                    alert(`Ошибка: ${error}`);
                }
            } catch (err) {
                console.error('Ошибка при получении пользователей:', err);
                alert('Произошла ошибка при получении пользователей');
            }
        });
    }
});

const apiBase = "http://localhost:8080";
let authToken = localStorage.getItem("authToken"); // Храните токен в localStorage после авторизации

// Загрузка профиля пользователя
async function loadUserProfile() {
    try {
        const response = await fetch(`${apiBase}/profile`, {
            headers: {
                "Authorization": `Bearer ${authToken}`,
            },
        });

        console.log('Profile response:', response); // Для проверки

        if (!response.ok) {
            throw new Error("Failed to load profile");
        }

        const user = await response.json();
        const profileName = document.getElementById("profile-name");
        const profileEmail = document.getElementById("profile-email");
        const profileImage = document.getElementById("profile-image");

        if (profileName) profileName.value = user.name || "";
        if (profileEmail) profileEmail.value = user.email || "";
        if (profileImage) profileImage.src = user.profile_picture || "default-profile.png";
        showProfileSection(user); // Показываем секцию профиля
    } catch (error) {
        console.error(error.message);
        alert("Please log in to access your profile.");
    }
}

// Обновление профиля
const profileForm = document.getElementById("profile-form");
if (profileForm) {
    profileForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const name = document.getElementById("profile-name")?.value;
        const email = document.getElementById("profile-email")?.value;

        try {
            const response = await fetch(`${apiBase}/profile/update`, {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${authToken}`,
                },
                body: JSON.stringify({ name, email }),
            });

            console.log('Profile update response:', response); // Для проверки

            if (!response.ok) {
                throw new Error("Failed to update profile");
            }

            alert("Profile updated successfully");
        } catch (error) {
            alert(error.message);
        }
    });
}

// Смена пароля
const passwordForm = document.getElementById("password-form");
if (passwordForm) {
    passwordForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const oldPassword = document.getElementById("old-password")?.value;
        const newPassword = document.getElementById("new-password")?.value;

        try {
            const response = await fetch(`${apiBase}/profile/password`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${authToken}`,
                },
                body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }),
            });

            console.log('Password change response:', response); // Для проверки

            if (!response.ok) {
                throw new Error("Failed to change password");
            }

            alert("Password changed successfully");
        } catch (error) {
            alert(error.message);
        }
    });
}

// Загрузка фотографии профиля
const pictureForm = document.getElementById("picture-form");
if (pictureForm) {
    pictureForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const fileInput = document.getElementById("profile-picture");
        const file = fileInput?.files[0];

        if (!file) {
            alert("Please select a file");
            return;
        }

        const formData = new FormData();
        formData.append("profile_picture", file);

        try {
            const response = await fetch(`${apiBase}/profile/picture`, {
                method: "POST",
                headers: {
                    "Authorization": `Bearer ${authToken}`,
                },
                body: formData,
            });

            console.log('Profile picture upload response:', response); // Для проверки

            if (!response.ok) {
                throw new Error("Failed to upload picture");
            }

            alert("Profile picture uploaded successfully");
            loadUserProfile();
        } catch (error) {
            alert(error.message);
        }
    });
}

// Загружаем профиль при открытии страницы
if (document.getElementById("profile-form")) {
    loadUserProfile();
}










