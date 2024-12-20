// Счетчики для ID элементов
let habitIdCounter = 1;

// Создание карточки элемента
function createCard(id, title, description, date, type, extra = '', actions = []) {
    const container = document.getElementById(`${type}-container`);
    const card = document.createElement('div');
    card.classList.add('card');
    card.dataset.id = id || 'No ID';

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



// Работа с привычками
// Работа с привычками

// Добавление привычки
async function addHabit() {
    const name = prompt('Enter habit name:');
    const description = prompt('Enter habit description:');
    const habit = { user_id: 1, name, description };

    const response = await fetch('http://localhost:8080/api/habits', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(habit)
    });

    if (response.ok) {
        const newHabit = await response.json();
        console.log('Habit added:', newHabit); // Для проверки
        getHabits(); // Обновление списка привычек
    } else {
        alert('Error adding habit.');
    }
}

// Получение привычек
async function getHabits() {
    const response = await fetch('http://localhost:8080/api/habits', { method: 'GET' });
    if (response.ok) {
        const habits = await response.json();
        const container = document.getElementById('habits-container');
        container.innerHTML = ''; // Очистка старых карточек

        habits.forEach(habit => {
            createCard(
                habit.id,
                habit.name,
                habit.description,
                new Date(habit.created_at).toLocaleString(),
                'habits',
                '',
                [
                    createActionButton('Edit', 'edit', () => editHabit(habit.id)),
                    createActionButton('Delete', 'delete', () => deleteHabit(habit.id))
                ]
            );
        });
    } else {
        alert('Error fetching habits.');
    }
}


async function editHabit(id) {
    const name = prompt('Enter new habit name:');
    const description = prompt('Enter new habit description:');

    const updatedHabit = { name, description };

    const response = await fetch(`http://localhost:8080/api/habits/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updatedHabit)
    });

    if (response.ok) {
        alert('Habit successfully updated!');
        getHabits(); // Обновление списка привычек
    } else {
        alert('Error updating habit.');
    }
}


// Удаление привычки
async function deleteHabit(id) {
    const response = await fetch(`http://localhost:8080/api/habits/${id}`, { method: 'DELETE' });
    if (response.ok) {
        alert('Habit successfully deleted!');
        getHabits(); // Обновление списка привычек
    } else {
        alert('Error deleting habit.');
    }
}


// Вызов получения привычек при загрузке страницы
window.onload = getHabits;


// Добавление цели
async function addGoal() {
    const name = prompt('Enter goal name:');
    const description = prompt('Enter goal description:');
    const deadline = prompt('Enter deadline (YYYY-MM-DD):');
    const goal = { user_id: 1, name, description, deadline };


    const response = await fetch('http://localhost:8080/api/goals', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(goal)
    });

    if (response.ok) {
        const newGoal = await response.json();
        console.log('Goal added:', newGoal); // Для проверки
        getGoals();

    } else {
        alert('Error adding goal.');
    }
}

// Изменение цели
async function editGoal(event) {
    const card = event.target.closest('.card');
    const name = card.querySelector('h3').textContent; // Получаем название цели
    const description = card.querySelector('p:first-of-type').textContent;
    const extra = card.querySelector('p:nth-of-type(4)');

    const newDescription = prompt('Edit goal description:', description);
    const newDeadline = prompt('Edit deadline (YYYY-MM-DD):', extra.textContent.split(': ')[1]);

    const updatedGoal = {
        name: name, // Название цели вместо ID
        description: newDescription || description,
        deadline: newDeadline || extra.textContent.split(': ')[1],
    };

    const response = await fetch('http://localhost:8080/api/goals/update-by-name', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(updatedGoal),
    });

    if (response.ok) {
        alert('Goal successfully updated!');
        card.querySelector('p:first-of-type').textContent = updatedGoal.description;
        extra.innerHTML = `<strong>Deadline:</strong> ${updatedGoal.deadline}`;
    } else {
        const error = await response.text();
        console.error('Error updating goal:', error);
        alert('Error updating goal.');
    }
}
async function getGoals() {
    const response = await fetch('http://localhost:8080/api/get-goals', { method: 'GET' });
    if (response.ok) {
        const goals = await response.json();
        const container = document.getElementById('goals-container');
        container.innerHTML = ''; // Очистка контейнера

        goals.forEach(goal => {
            createCard(
                goal.id,
                goal.name,
                goal.description,
                new Date(goal.created_at).toLocaleString(),
                "goals",
                `<strong>Deadline:</strong> ${goal.deadline}`,
                [
                    createActionButton('Edit', 'edit', () => editGoal(goal.name)),
                    createActionButton('Delete', 'delete', () => deleteGoal(goal.id))
                ]
            );
        });
    } else {
        console.error('Error fetching goals:', await response.text());
        alert('Error fetching goals.');
    }
}



// Удаление цели
async function deleteGoal(event) {
    const card = event.target.closest('.card');
    const id = card.dataset.id;

    const response = await fetch(`http://localhost:8080/api/goals/${id}`, {
        method: 'DELETE'
    });

    if (response.ok) {
        card.remove();
        alert('Goal successfully deleted!');
    } else {
        alert('Error deleting goal.');
    }
}
window.onload = getGoals();

// Добавление достижения
// Работа с достижениями (Achievements)

// Добавление достижения
async function addAchievement() {
    const title = prompt('Enter achievement title:');
    const description = prompt('Enter achievement description:');
    const achievement = {
        user_id: 1, // ID пользователя
        title,
        description,
        date: new Date().toISOString() // Текущая дата
    };

    const response = await fetch('http://localhost:8080/api/achievements', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(achievement)
    });

    if (response.ok) {
        const newAchievement = await response.json();
        console.log('Achievement added:', newAchievement);
        getAchievements(); // Обновление списка достижений
    } else {
        console.error('Error adding achievement:', await response.text());
        alert('Error adding achievement.');
    }
}

// Получение списка достижений
async function getAchievements() {
    const response = await fetch('http://localhost:8080/api/achievements', { method: 'GET' });
    if (response.ok) {
        const achievements = await response.json();
        const container = document.getElementById('achievements-container');
        container.innerHTML = ''; // Очистка старых карточек

        achievements.forEach(achievement => {
            createCard(
                achievement.id,
                achievement.title,
                achievement.description,
                new Date(achievement.date).toLocaleString(),
                'achievements',
                '',
                [
                    createActionButton('Delete', 'delete', () => deleteAchievement(achievement.id))
                ]
            );
        });
    } else {
        console.error('Error fetching achievements:', await response.text());
        alert('Error fetching achievements.');
    }
}

// Удаление достижения
async function deleteAchievement(id) {
    const response = await fetch(`http://localhost:8080/api/achievements/${id}`, { method: 'DELETE' });
    if (response.ok) {
        alert('Achievement successfully deleted!');
        getAchievements(); // Обновление списка достижений
    } else {
        console.error('Error deleting achievement:', await response.text());
        alert('Error deleting achievement.');
    }
}

// Вызов получения достижений при загрузке страницы
window.onload = function () {
    getAchievements();
};


// Добавление уведомления
async function addNotification() {
    const message = prompt('Enter notification message:'); // Запрос сообщения
    const notification = {
        user_id: 1,
        message,
        scheduled_at: new Date().toISOString() // Планируемая дата
    };

    const response = await fetch('http://localhost:8080/api/notifications', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(notification)
    });

    if (response.ok) {
        const newNotification = await response.json();
        createCard(
            newNotification.id,
            'Notification',
            newNotification.message,
            new Date(newNotification.scheduled_at).toLocaleString(),
            'notifications',
            `<strong>Is Sent:</strong> ${newNotification.is_sent ? 'Yes' : 'No'}`,
            [
                createActionButton('Delete', 'delete', () => deleteNotification(newNotification.id))
            ]
        );
        alert('Notification successfully added!');
    } else {
        console.error('Error adding notification:', await response.text());
        alert('Error adding notification.');
    }
}

// Получение списка уведомлений
async function getNotifications() {
    const response = await fetch('http://localhost:8080/api/notifications', { method: 'GET' });
    if (response.ok) {
        const notifications = await response.json();
        const container = document.getElementById('notifications-container');
        container.innerHTML = ''; // Очистка контейнера

        notifications.forEach(notification => {
            createCard(
                notification.id,
                'Notification',
                notification.message,
                new Date(notification.scheduled_at).toLocaleString(),
                'notifications',
                `<strong>Is Sent:</strong> ${notification.is_sent ? 'Yes' : 'No'}`,
                [
                    createActionButton('Delete', 'delete', () => deleteNotification(notification.id))
                ]
            );
        });
    } else {
        console.error('Error fetching notifications:', await response.text());
        alert('Error fetching notifications.');
    }
}

// Удаление уведомления
async function deleteNotification(id) {
    const response = await fetch(`http://localhost:8080/api/notifications/${id}`, { method: 'DELETE' });
    if (response.ok) {
        alert('Notification successfully deleted!');
        getNotifications(); // Обновить список уведомлений
    } else {
        console.error('Error deleting notification:', await response.text());
        alert('Error deleting notification.');
    }
}


window.onload = function () {
    getHabits();
    getGoals();
    getAchievements();
    getNotifications();
};
