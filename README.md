**Step-by-Step Detailed Explanation:**

1. **Context:**  
   HabitMaster is a web-based habit tracking and productivity application. It allows users to create, monitor, and maintain habits, visualize their progress, and receive reminders or notifications. The project is written in Go and uses Gorilla/mux for HTTP routing.

2. **Purpose of README:**  
   This README offers a clear overview of the project, highlighting its main functionalities, intended audience, instructions for local setup and launch, as well as introducing the development team.

3. **Project Setup and Running Instructions:**  
   - **Run the server:** After configuring and running your Go application, you can access the HabitMaster interface via `http://localhost:8080/main.html` in your browser.
   - **Stop the server:** Press `Ctrl + C` in the terminal where the server is running to stop the project.
   
   **Code snippet example (without placeholders):**
   ```go
   import (
       "HabitMaster/databaseConnector"
       "HabitMaster/handlers"
       "fmt"
       "github.com/gorilla/mux"
       "log"
       "net/http"
   )

   func main() {
       r := mux.NewRouter()
       // Sample route
       r.HandleFunc("/main.html", handlers.MainHandler)

       fmt.Println("Server running at http://localhost:8080/main.html")
       if err := http.ListenAndServe(":8080", r); err != nil {
           log.Fatal(err)
       }
   }
   ```

4. **Key Features:**
   - **Custom Habit Creation:** Set daily, weekly, or monthly habits with fully customizable goals and schedules.  
   - **Visual Progress Tracking:** View interactive charts, streaks, and completion rates, helping you understand your performance over time.  
   - **Reminders and Notifications:** Get email or push notifications to ensure you stay on track with your habits.

5. **Audience:**
   Whether you’re a student aiming to develop better study routines, a professional working on productivity goals, or an individual seeking personal growth, HabitMaster is designed to guide you toward building and maintaining positive habits.

6. **Team Members:**
   - Zayniddinov Berik  
   - Rishat Nurassyl  
   - Myrzan Myrzakhan

7. **Screenshot of the Main Page:**
   ![image](https://github.com/user-attachments/assets/75d7f1eb-ca6c-4246-951c-fbdbaa137913)


---

**Conclusion:**  
HabitMaster is more than just a habit tracker—it’s your digital accountability partner. By following the instructions above, you can quickly run the project locally, explore its features, and begin transforming your daily routines into long-term successes.





**2 Assignment**

```markdown
# HabitMaster

Welcome to **HabitMaster**, an advanced habit and goal tracking system designed to improve productivity and simplify the management of your daily tasks. This project is built using modern web technologies, including **Go (Golang)** for the backend, **PostgreSQL** for database management, and a responsive user interface powered by **HTML, CSS, and JavaScript**.

This repository demonstrates a variety of advanced web development concepts, such as **Filtering, Sorting, Pagination**, **Structured Logging**, **Error Handling**, **Rate Limiting**, **Graceful Shutdown**, and **Email Integration**. Below, you'll find an overview of the project's key features and implementation details.

---

## Features

### Core Functionalities:
1. **Filtering, Sorting, and Pagination**:
   - Dynamic filtering and sorting of habits, goals, achievements, and notifications.
   - Pagination implemented on both backend and frontend to optimize performance.

2. **Structured Logging**:
   - Logs user actions, such as filtering, sorting, and data manipulation.
   - Uses `logrus` for structured logging in JSON format.

3. **Error Handling**:
   - Robust error handling mechanisms to ensure smooth user experience.
   - Displays descriptive error messages for invalid inputs or server issues.

4. **Rate Limiting**:
   - Prevents server overload by limiting the number of requests per client using `golang.org/x/time/rate`.

5. **Graceful Shutdown**:
   - Ensures the server handles termination signals properly and completes all active requests before shutting down.

6. **Email Sending**:
   - Administrators can send promotional or informational emails to users.
   - Users can send support requests via email directly from the application.
   - Integrated with SMTP using `gomail.v2`.

---

## Setup Instructions

### Prerequisites:
- **Golang** installed on your machine.
- **PostgreSQL** database configured.
- A valid SMTP account for email services.

### Steps:
1. Clone this repository:
   ```bash
   git clone https://github.com/berikzayniddinov/HabitMaster.git
   cd HabitMaster
   ```

2. Create a `.env` file in the root directory and configure the following variables:
   ```
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_database_user
   DB_PASSWORD=your_database_password
   DB_NAME=your_database_name

   SMTP_HOST=smtp.gmail.com
   SMTP_PORT=587
   SMTP_USER=your_email@gmail.com
   SMTP_PASSWORD=your_email_password
   ```

3. Run the application:
   ```bash
   go run main/main.go
   ```

4. Access the app at `http://localhost:8080` in your browser.

---

## Technologies Used

- **Backend**: Go (Golang)
- **Database**: PostgreSQL
- **Frontend**: HTML, CSS, JavaScript
- **Libraries**:
  - `github.com/gorilla/mux` (Routing)
  - `github.com/sirupsen/logrus` (Logging)
  - `gopkg.in/gomail.v2` (Email Sending)
  - `golang.org/x/time/rate` (Rate Limiting)

---

## Project Structure

- **`main/`**: Entry point for the application.
- **`handlers/`**: Contains API endpoint handlers for habits, goals, achievements, notifications, and emails.
- **`habittracker/`**: Frontend files including HTML, CSS, and JavaScript.
- **`emailSender/`**: Logic for sending emails using SMTP.
- **`databaseConnector/`**: Database connection setup and queries.

---

## Key Features in Detail

### Filtering, Sorting, and Pagination:
- Implemented dynamic query building in Go to support filtering and sorting parameters.
- Pagination helps load data in smaller chunks, reducing server load and improving user experience.

### Structured Logging:
- All actions and errors are logged in JSON format using `logrus`.
- Example log:
  ```json
  {
    "time": "2025-01-09T12:00:00Z",
    "level": "info",
    "action": "fetch_data",
    "status": "success",
    "module": "notifications",
    "message": "Notifications fetched successfully"
  }
  ```

### Email Integration:
- Admins can send promotional emails to multiple users.
- Users can send support requests from the profile page.

---

## Contributing

Feel free to fork this repository, make your changes, and submit a pull request. Contributions are always welcome!

---

## License

This project is licensed under the [MIT License](LICENSE).
```

Эта версия обновленного `README.md` начинается с более формального и профессионального описания проекта. Она охватывает все основные задания из вашего ассаймента, а также добавляет дополнительные секции для удобства разработчиков и пользователей.

```

3 assugnment soon
