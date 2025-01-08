# Forum Project

## Objectives

This project is a web forum with the following functionality:

- **Communication between users**: Users can interact by creating posts and comments.
- **Categorization of posts**: Posts can be associated with one or more categories.
- **Likes and dislikes**: Users can like or dislike posts and comments, with the counts visible to everyone.
- **Filtering posts**: Posts can be filtered by categories, user-created posts, and liked posts.

---

## SQLite

SQLite is used to store the forum's data (e.g., users, posts, comments). It is an embedded database software ideal for local storage in application software.

### Notes:

SQLite enables creating and controlling a database using queries. To learn more about SQLite, visit the [SQLite documentation](https://sqlite.org/).

---

## Authentication

The forum supports user authentication, including:

### Registration:

- **Email**:
  - Collecting an email address during registration.
  - If the email is already in use, error response is returned.
- **Username**:
  - Collecting a unique username.
- **Password**:
  - Collecting a password and optionally encrypt it when stored.

### Login:

- The forum verifies if the provided email exists in the database and check if the credentials are correct.
- If the password doesn't match, it returns an error response.

### Sessions:

- Use of **cookies** to manage user sessions.

---

## Communication

To facilitate communication among users:

- **Registered users**:
  - Can create posts and comments.
  - Posts can be associated with one or more categories (you decide the categories).
- **Non-registered users**:
  - Can only view posts and comments.

---

## Likes and Dislikes

- Only registered users can like or dislike posts and comments.
- The total number of likes and dislikes is visible to everyone (registered or not).

---

## Filter

The forum includes a filtering mechanism to:

- Filter posts by **categories** (like subforums for specific topics).
- Display posts created by the logged-in user (**created posts**).
- Display posts liked by the logged-in user (**liked posts**).

### Notes:

- The "created posts" and "liked posts" filters is only available to registered users.

---

## Docker

**Docker** has been used to allow packaging the application and its dependencies into a container, ensuring consistent behavior across environments.

To build the image

```
docker build -t forum .
```

Then to run the bult image

```
docker run -d -p 8000:8000 forum
```

---

## How to run the application

1. Clone the Repository
   
   ```
   git clone https://learn.zone01kisumu.ke/git/aaochieng/forum.git

   cd forum
   ```
2. run the following command
    ```
    make
    ```
    or

    ```
    go run ./cmd/web/
    ```

3. On your Web Browser,
   
   ```
   localhost:8000
   ```
4. To run tests:
   
   ```
   make test
   ```
   
   or
   
   ```
   go test ./...
   ```

## Authors

[Aaron Ochieng](https://github.com/Aaron-Ochieng)

[Moses Onyango](https://github.com/moseeh)

[Swabri Musa](https://github.com/skanenje)

[Andrew Osindo](https://github.com/andyosyndoh)