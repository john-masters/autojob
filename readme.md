# getjob.xyz

## IDEA

- A site where you create an account and add your resume and cover letter.
- It runs a cron job daily that scrapes job site listings after searching with your preferred search term.
- It then goes through each job and uses GPT-4 to decide if the job is a good fit for you.
- If it is a good fit, it creates a custom cover letter using GPT-4 and adds it to a `To Apply` list.
  - Ideally it would apply to the job for you, but job posts rarely have an email address to apply to and often have custom application forms.

## TECH

- Built with Go
  - [Templ](https://templ.guide/) for html templating
  - [Colly](https://go-colly.org/) for web scraping
- SQLite for the database
- BCrypt for password hashing
- JWT for authentication

## COMMANDS

- Run the server: `go run .`
- Update the html: `templ generate`
- Watch and regenerate the html:

  ```sh
  templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."
  ```

- Build the binary: `go build .`

## SETUP

1. Install dependencies: `go mod tidy`
1. Add environment variables
    1. Rename `.env.sample` to `.env`
    1. Add your environment variables
        - To genererate a JWT secret, run `openssl rand -hex 32` or just use a random string
1. Run the server: `go run .`
    - Optional: uncomment the `// utils.DbInit()` line in `main.go` to initialise the db
1. Visit `http://localhost:8080`
1. Set admin:
    - Connect to db using CLI: `sqlite3 test.db`
    - Update permission: `UPDATE users SET is_admin = 1 WHERE id = <your user_id>;`
    - Quit: `.quit`

## TODO

- [x] Add user auth
- [x] Add account page
- [x] Add settings page
- [x] Add user resume and cover letter upload
- [x] Add web scraping
- [x] Add GPT-4 integration
- [x] Add cron job to run daily
- [x] Add more robust input validation
  - [x] valid email
  - [x] password complexity
- [ ] Fix job description scraping bug (no new lines)
- [ ] Add indeed to scraped sites

