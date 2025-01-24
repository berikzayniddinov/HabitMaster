// admin.js

const apiBaseUrl = 'https://abc123.ngrok.io'; // Замените на актуальный URL ngrok

// Общая функция для выполнения запросов
async function fetchData(endpoint, options = {}) {
    try {
        const response = await fetch(`${apiBaseUrl}${endpoint}`, options);
        if (!response.ok) {
            throw new Error(`Error: ${response.statusText}`);
        }
        return await response.json();
    } catch (error) {
        console.error('Fetch error:', error);
        alert('An error occurred. Please check the console for details.');
        return null;
    }
}
async function sendPromotionalEmail() {
    const to = document.getElementById('to').value;
    const subject = document.getElementById('subject').value;
    const body = document.getElementById('body').value;
    const attachment = document.getElementById('attachment').files[0];

    // Проверка на заполненность обязательных полей
    if (!to || !subject || !body) {
        alert('Please fill out all required fields!');
        return;
    }

    // Формирование данных для отправки
    const formData = new FormData();
    formData.append('to', to);
    formData.append('subject', subject);
    formData.append('body', body);
    if (attachment) {
        formData.append('attachment', attachment);
    }

    // Отправка данных на сервер
    try {
        const response = await fetch(`${apiBaseUrl}/api/admin/send-mass-email`, {
            method: 'POST',
            body: formData, // Использование FormData для загрузки файла
        });

        if (response.ok) {
            alert('Email sent successfully!');
            document.getElementById('email-form').reset();
        } else {
            const errorText = await response.text();
            throw new Error(errorText);
        }
    } catch (error) {
        console.error('Error sending email:', error);
        alert('Failed to send email. Please check the console for details.');
    }
}

async function getUsers() {
    const users = await fetchData('/api/get-users');
    if (users) {
        const userList = document.getElementById('user-list');
        userList.innerHTML = '';
        users.forEach(user => {
            const li = document.createElement('li');
            li.textContent = `${user.name} - ${user.email}`;
            userList.appendChild(li);
        });
    }
}

// Применение фильтров для пользователей
async function applyUserFilters() {
    const emailFilter = document.getElementById('emailFilter').value;
    const usernameFilter = document.getElementById('usernameFilter').value;

    const query = new URLSearchParams();
    if (emailFilter) query.append('email', emailFilter);
    if (usernameFilter) query.append('username', usernameFilter);

    const users = await fetchData(`/api/get-users?${query.toString()}`);
    if (users) {
        const userList = document.getElementById('user-list');
        userList.innerHTML = '';
        users.forEach(user => {
            const li = document.createElement('li');
            li.textContent = `${user.name} - ${user.email}`;
            userList.appendChild(li);
        });
    }
}

// Применение сортировки пользователей
async function applyUserSort() {
    const sortOption = document.getElementById('sortUsers').value;

    const query = new URLSearchParams();
    query.append('sort', sortOption);

    const users = await fetchData(`/api/get-users?${query.toString()}`);
    if (users) {
        const userList = document.getElementById('user-list');
        userList.innerHTML = '';
        users.forEach(user => {
            const li = document.createElement('li');
            li.textContent = `${user.name} - ${user.email}`;
            userList.appendChild(li);
        });
    }
}

// Переход на определённую страницу пользователей
async function goToUserPage() {
    const page = document.getElementById('pageNumber').value;

    const query = new URLSearchParams();
    query.append('page', page);

    const users = await fetchData(`/api/get-users?${query.toString()}`);
    if (users) {
        const userList = document.getElementById('user-list');
        userList.innerHTML = '';
        users.forEach(user => {
            const li = document.createElement('li');
            li.textContent = `${user.name} - ${user.email}`;
            userList.appendChild(li);
        });
    }
}

// Создание нового пользователя
async function createUser() {
    const name = prompt('Enter user name:');
    const email = prompt('Enter user email:');
    const password = prompt('Enter user password:');

    if (!name || !email || !password) {
        alert('All fields are required!');
        return;
    }

    const user = { name, email, password };

    try {
        const response = await fetch(`${apiBaseUrl}/api/users`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(user),
        });

        if (response.ok) {
            alert('User created successfully!');
            getUsers();
        } else {
            throw new Error(await response.text());
        }
    } catch (error) {
        console.error('Error creating user:', error);
        alert('Failed to create user.');
    }
}

// Обновление существующего пользователя
async function updateUser() {
    const id = prompt('Enter user ID to update:');
    const name = prompt('Enter new user name:');
    const email = prompt('Enter new user email:');

    if (!id || !name || !email) {
        alert('ID, name, and email are required!');
        return;
    }

    const user = { id, name, email };

    try {
        const response = await fetch(`${apiBaseUrl}/api/users`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(user),
        });

        if (response.ok) {
            alert('User updated successfully!');
            getUsers();
        } else {
            throw new Error(await response.text());
        }
    } catch (error) {
        console.error('Error updating user:', error);
        alert('Failed to update user.');
    }
}

// Удаление пользователя
async function deleteUser() {
    const id = prompt('Enter user ID to delete:');

    if (!id) {
        alert('ID is required!');
        return;
    }

    try {
        const response = await fetch(`${apiBaseUrl}/api/users?id=${id}`, {
            method: 'DELETE',
        });

        if (response.ok) {
            alert('User deleted successfully!');
            getUsers();
        } else {
            throw new Error(await response.text());
        }
    } catch (error) {
        console.error('Error deleting user:', error);
        alert('Failed to delete user.');
    }
}

// Поиск пользователя по email
async function searchUserByEmail() {
    const email = document.getElementById('searchEmail').value;

    if (!email) {
        alert('Email is required!');
        return;
    }

    const user = await fetchData(`/api/get-users?email=${encodeURIComponent(email)}`);
    if (user && user.length > 0) {
        alert(`User found: ${user[0].name} - ${user[0].email}`);
    } else {
        alert('User not found.');
    }
}

// Поиск пользователя по имени пользователя
async function searchUserByUsername() {
    const username = document.getElementById('searchUsername').value;

    if (!username) {
        alert('Username is required!');
        return;
    }

    const users = await fetchData(`/api/get-users?username=${encodeURIComponent(username)}`);
    if (users && users.length > 0) {
        alert(`User(s) found: ${users.map(user => `${user.name} - ${user.email}`).join(', ')}`);
    } else {
        alert('User not found.');
    }
}


// Получение списка привычек
async function getHabits() {
    const habits = await fetchData('/api/habits');
    if (habits) {
        const habitList = document.getElementById('habit-list');
        habitList.innerHTML = '';
        habits.forEach(habit => {
            const li = document.createElement('li');
            li.textContent = `${habit.name} - ${habit.description}`;
            habitList.appendChild(li);
        });
    }
}

async function getGoals() {
    const goals = await fetchData('/api/goals');
    if (goals) {
        const goalList = document.getElementById('goal-list');
        goalList.innerHTML = '';
        goals.forEach(goal => {
            const li = document.createElement('li');
            li.textContent = `${goal.name} - ${goal.deadline}`;
            goalList.appendChild(li);
        });
    }
}

// Применение фильтров для целей
async function applyGoalFilters() {
    const nameFilter = document.getElementById('goal-filter-name').value;

    const query = new URLSearchParams();
    if (nameFilter) query.append('name', nameFilter);

    const goals = await fetchData(`/api/goals?${query.toString()}`);
    if (goals) {
        const goalList = document.getElementById('goal-list');
        goalList.innerHTML = '';
        goals.forEach(goal => {
            const li = document.createElement('li');
            li.textContent = `${goal.name} - ${goal.deadline}`;
            goalList.appendChild(li);
        });
    }
}

// Переход на определённую страницу целей
async function goToGoalPage() {
    const page = document.getElementById('goal-page').value;

    const query = new URLSearchParams();
    query.append('page', page);

    const goals = await fetchData(`/api/goals?${query.toString()}`);
    if (goals) {
        const goalList = document.getElementById('goal-list');
        goalList.innerHTML = '';
        goals.forEach(goal => {
            const li = document.createElement('li');
            li.textContent = `${goal.name} - ${goal.deadline}`;
            goalList.appendChild(li);
        });
    }
}

// Создание новой цели
async function createGoal() {
    const name = prompt('Enter goal name:');
    const description = prompt('Enter goal description:');
    const deadline = prompt('Enter goal deadline (YYYY-MM-DD):');

    if (!name || !deadline) {
        alert('Name and deadline are required!');
        return;
    }

    const goal = { name, description, deadline };

    try {
        const response = await fetch(`${apiBaseUrl}/api/goals`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(goal),
        });

        if (response.ok) {
            alert('Goal created successfully!');
            getGoals();
        } else {
            throw new Error(await response.text());
        }
    } catch (error) {
        console.error('Error creating goal:', error);
        alert('Failed to create goal.');
    }
}

// Обновление существующей цели
async function updateGoal() {
    const id = prompt('Enter goal ID to update:');
    const name = prompt('Enter new goal name:');
    const description = prompt('Enter new goal description:');
    const deadline = prompt('Enter new goal deadline (YYYY-MM-DD):');

    if (!id || !name || !deadline) {
        alert('ID, name, and deadline are required!');
        return;
    }

    const goal = { id, name, description, deadline };

    try {
        const response = await fetch(`${apiBaseUrl}/api/goals`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(goal),
        });

        if (response.ok) {
            alert('Goal updated successfully!');
            getGoals();
        } else {
            throw new Error(await response.text());
        }
    } catch (error) {
        console.error('Error updating goal:', error);
        alert('Failed to update goal.');
    }
}

// Удаление цели
async function deleteGoal() {
    const id = prompt('Enter goal ID to delete:');

    if (!id) {
        alert('ID is required!');
        return;
    }

    try {
        const response = await fetch(`${apiBaseUrl}/api/goals?id=${id}`, {
            method: 'DELETE',
        });

        if (response.ok) {
            alert('Goal deleted successfully!');
            getGoals();
        } else {
            throw new Error(await response.text());
        }
    } catch (error) {
        console.error('Error deleting goal:', error);
        alert('Failed to delete goal.');
    }
}

// Поиск цели по ID
async function searchGoalById() {
    const id = document.getElementById('goal-id').value;

    if (!id) {
        alert('ID is required!');
        return;
    }

    const goal = await fetchData(`/api/goals?id=${id}`);
    if (goal) {
        alert(`Goal found: ${goal.name} - ${goal.description}`);
    } else {
        alert('Goal not found.');
    }
}

// Поиск цели по имени
async function searchGoalByName() {
    const name = document.getElementById('goal-name').value;

    if (!name) {
        alert('Name is required!');
        return;
    }

    const goals = await fetchData(`/api/goals?name=${encodeURIComponent(name)}`);
    if (goals && goals.length > 0) {
        const goalList = document.getElementById('goal-list');
        goalList.innerHTML = '';
        goals.forEach(goal => {
            const li = document.createElement('li');
            li.textContent = `${goal.name} - ${goal.deadline}`;
            goalList.appendChild(li);
        });
    } else {
        alert('No goals found.');
    }
}
// Получение списка привычек
async function getHabits() {
    const habits = await fetchData('/api/habits');
    if (habits) {
        const habitList = document.getElementById('habit-list');
        habitList.innerHTML = '';
        habits.forEach(habit => {
            const li = document.createElement('li');
            li.textContent = `${habit.name} - ${habit.description}`;
            habitList.appendChild(li);
        });
    }
}

// Применение фильтров для привычек
async function applyHabitFilters() {
    const nameFilter = document.getElementById('habit-filter-name').value;

    const query = new URLSearchParams();
    if (nameFilter) query.append('name', nameFilter);

    const habits = await fetchData(`/api/habits?${query.toString()}`);
    if (habits) {
        const habitList = document.getElementById('habit-list');
        habitList.innerHTML = '';
        habits.forEach(habit => {
            const li = document.createElement('li');
            li.textContent = `${habit.name} - ${habit.description}`;
            habitList.appendChild(li);
        });
    }
}

// Переход на определённую страницу привычек
async function goToHabitPage() {
    const page = document.getElementById('habit-page').value;

    const query = new URLSearchParams();
    query.append('page', page);

    const habits = await fetchData(`/api/habits?${query.toString()}`);
    if (habits) {
        const habitList = document.getElementById('habit-list');
        habitList.innerHTML = '';
        habits.forEach(habit => {
            const li = document.createElement('li');
            li.textContent = `${habit.name} - ${habit.description}`;
            habitList.appendChild(li);
        });
    }
}

// Создание новой привычки
async function createHabit() {
    const name = prompt('Enter habit name:');
    const description = prompt('Enter habit description:');

    if (!name) {
        alert('Name is required!');
        return;
    }

    const habit = { name, description };

    try {
        const response = await fetch(`${apiBaseUrl}/api/habits`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(habit),
        });

        if (response.ok) {
            alert('Habit created successfully!');
            getHabits();
        } else {
            throw new Error(await response.text());
        }
    } catch (error) {
        console.error('Error creating habit:', error);
        alert('Failed to create habit.');
    }
}

// Обновление существующей привычки
async function updateHabit() {
    const id = prompt('Enter habit ID to update:');
    const name = prompt('Enter new habit name:');
    const description = prompt('Enter new habit description:');

    if (!id || !name) {
        alert('ID and name are required!');
        return;
    }

    const habit = { id, name, description };

    try {
        const response = await fetch(`${apiBaseUrl}/api/habits`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(habit),
        });

        if (response.ok) {
            alert('Habit updated successfully!');
            getHabits();
        } else {
            throw new Error(await response.text());
        }
    } catch (error) {
        console.error('Error updating habit:', error);
        alert('Failed to update habit.');
    }
}

// Удаление привычки
async function deleteHabit() {
    const id = prompt('Enter habit ID to delete:');

    if (!id) {
        alert('ID is required!');
        return;
    }

    try {
        const response = await fetch(`${apiBaseUrl}/api/habits?id=${id}`, {
            method: 'DELETE',
        });

        if (response.ok) {
            alert('Habit deleted successfully!');
            getHabits();
        } else {
            throw new Error(await response.text());
        }
    } catch (error) {
        console.error('Error deleting habit:', error);
        alert('Failed to delete habit.');
    }
}

// Поиск привычки по ID
async function searchHabitById() {
    const id = document.getElementById('habit-id').value;

    if (!id) {
        alert('ID is required!');
        return;
    }

    const habit = await fetchData(`/api/habits?id=${id}`);
    if (habit) {
        alert(`Habit found: ${habit.name} - ${habit.description}`);
    } else {
        alert('Habit not found.');
    }
}

// Поиск привычки по имени
async function searchHabitByName() {
    const name = document.getElementById('habit-name').value;

    if (!name) {
        alert('Name is required!');
        return;
    }

    const habits = await fetchData(`/api/habits?name=${encodeURIComponent(name)}`);
    if (habits && habits.length > 0) {
        const habitList = document.getElementById('habit-list');
        habitList.innerHTML = '';
        habits.forEach(habit => {
            const li = document.createElement('li');
            li.textContent = `${habit.name} - ${habit.description}`;
            habitList.appendChild(li);
        });
    } else {
        alert('No habits found.');
    }
}