package email

// import (
// 	"fmt"
// 	ticketModel "ideyanale-be/pkg/modules/tickets/model"
// 	"log"
// 	"net/smtp"
// 	"os"
// )

// // SendEndorserNotification — called on ticket creation
// func SendEndorserNotification(ticket ticketModel.Ticket, toEmail string, submitterName string) error {
// 	from := os.Getenv("EMAIL_ADDRESS")
// 	password := os.Getenv("EMAIL_PASSWORD")
// 	smtpHost := os.Getenv("SMTP_HOST")

// 	auth := smtp.PlainAuth("", from, password, smtpHost)

// 	subject := "New Ticket For Endorsement"

// 	body := fmt.Sprintf(`
// <!DOCTYPE html>
// <html>
// <head>
//   <style>
//     body {
//       font-family: Arial, sans-serif;
//       background-color: #f4f6f8;
//       margin: 0;
//       padding: 20px;
//     }
//     .container {
//       max-width: 600px;
//       margin: auto;
//       background: #ffffff;
//       padding: 25px;
//       border-radius: 10px;
//       box-shadow: 0 4px 10px rgba(0,0,0,0.08);
//     }
//     .header {
//       text-align: center;
//       padding-bottom: 10px;
//       border-bottom: 1px solid #eee;
//     }
//     .header h2 {
//       color: #2c3e50;
//       margin: 0;
//     }
//     .ticket-box {
//       background: #f9fafb;
//       padding: 15px;
//       margin-top: 20px;
//       border-radius: 8px;
//       border-left: 5px solid #3498db;
//     }
//     .label {
//       font-weight: bold;
//       color: #555;
//     }
//     .value {
//       color: #222;
//     }
//     .btn {
//   display: block;
//   text-align: center;
//   margin: 25px auto 10px;
//   padding: 12px;
//   font-size: 16px;
//   color: #ffffff !important;
//   background-color: #007bff;
//   text-decoration: none !important;
//   border-radius: 6px;
//   width: 200px;
// }
//     .btn:hover {
//       background-color: #0056b3;
//     }
//     .footer {
//       margin-top: 25px;
//       font-size: 12px;
//       text-align: center;
//       color: #888;
//       border-top: 1px solid #eee;
//       padding-top: 15px;
//     }
//     .note {
//       font-size: 12px;
//       color: #999;
//       margin-top: 5px;
//     }
//   </style>
// </head>
// <body>
//   <div class="container">
    
//     <div class="header">
//       <h2>New Ticket For Endorsement</h2>
//       <p>You have received a ticket that requires your action.</p>
//     </div>

//     <div class="ticket-box">
//       <p><span class="label">Ticket ID:</span> %s</p>
//       <p><span class="label">Subject:</span> %s</p>
//       <p><span class="label">Category:</span> %s</p>
//       <p><span class="label">Priority:</span> %s</p>
//       <p><span class="label">Submitted By:</span> %s</p>
//     </div>

//     <!-- LOGIN BUTTON -->
//     <a href="https://idiyanale.bakawan-ai.com/login" class="btn">Go to Login</a>

//     <div class="footer">
//       <p><b>Note:</b> This message is auto-generated.</p>
//       <p>Please do not reply to this email.</p>
//       <div class="note">If you have concerns, contact the system administrator.</div>
//     </div>

//   </div>
// </body>
// </html>
// `, ticket.TicketID, ticket.Subject, ticket.Category, ticket.Priority, submitterName)

// 	msg := []byte(
// 		"Subject: " + subject + "\r\n" +
// 			"MIME-Version: 1.0\r\n" +
// 			"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
// 			body,
// 	)

// 	err := smtp.SendMail(smtpHost+":587", auth, from, []string{toEmail}, msg)
// 	if err != nil {
// 		log.Println("Failed to send endorser email:", err)
// 	}

// 	return err
// }
