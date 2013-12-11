gomail
======

Render plain text or html email body from Go template and send it. It supports TLS and written in Go.

Here is an example:

```Go
package main

import (
  "fmt"

  "github.com/arkxu/gomail"
)

func main() {
  // init the mailer and tell the smtp server information
  mailer := &gomail.Mailer{Server: "smtp.gmail.com", Port: 587, UserName: "your gmail", Password: "your pwd"}

  // set the default sender. this is optional(you can set it in the message itself)
  mailer.Sender = &gomail.Sender{From: "foo@bar.com"}

  // configure the template folder, the file with extension .txt will be rendered as plain text email body
  // and the file with extension .html will be rendered as html email body
  err := gomail.TemplateFolder("templates")
  if err != nil {
    panic(err)
  }

  // arguments used for template rendering
  var args = make(map[string]interface{})
  args["world"] = "世界"
  args["user"] = struct {
    Name string
    Link string
  }{
    "Go",
    "http://golang.org",
  }

  // init the message, set the CC or BCC, Subject or body
  message := &gomail.Message{To: []string{"anyone@foo.com"}, Subject: "from template 6", Cc: []string{"bar@foo.com"}}

  // set the template file for this message
  err = message.RenderTemplate("helloWorld", args)
  if err != nil {
    fmt.Println(err)
  }

  // send the message out
  err = mailer.SendMessage(message)
  if err != nil {
    fmt.Println(err)
  }
}

```

In order to make it work, you need to create a folder `templates` and put the template file in there.

`TemplateFolder()` is not mandatory if the email body is passed in, instead of rendered. Please check out the test cases
for more information.

