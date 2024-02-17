package main

////func listenForMail() {
////	go func() {
////		for {
////			m := <-app.MailChan
////			sendMsg(m)
////		}
////	}()
////}
//
//func sendMsg(m models.MailData) {
//	server := mail.NewSMTPClient()
//	server.Host = "localhost"
//	server.Port = 1025
//	server.KeepAlive = false
//	server.ConnectTimeout = 10 * time.Second
//	server.SendTimeout = 10 * time.Second
//
//	client, err := server.Connect()
//	if err != nil {
//		app.ErrorLog.Println("error connecting mail server:", err)
//	}
//	email := mail.NewMSG()
//	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
//	if m.Template == "" {
//		mailContent, err := os.ReadFile(fmt.Sprintf("./email-templates/%s", m.Content))
//		if err != nil {
//			app.ErrorLog.Println("error read email template")
//		}
//		msgToSend := string(mailContent)
//		email.SetBody(mail.TextHTML, msgToSend)
//	} else {
//		mailTemplate, err := os.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
//		if err != nil {
//			app.ErrorLog.Println("error read email template")
//		}
//		mailContent, err := os.ReadFile(fmt.Sprintf("./email-templates/%s", m.Content))
//		if err != nil {
//			app.ErrorLog.Println("error read email template")
//		}
//		msgToSend := strings.Replace(string(mailTemplate), "[%body%]", string(mailContent), 1)
//
//		email.SetBody(mail.TextHTML, msgToSend)
//	}
//
//	err = email.Send(client)
//	if err != nil {
//		app.ErrorLog.Println("error sending email:", err)
//	} else {
//		app.InfoLog.Println("email sent")
//	}
//
//}
