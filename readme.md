# AUTOJOB

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
2. Run the server: `go run .`
    - Optional: uncomment the `// utils.DbInit()` line in `main.go` to initialise the db
3. Visit `http://localhost:8080`

## TODO

- Add user auth
- Add user resume and cover letter upload
