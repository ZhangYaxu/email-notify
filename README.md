# Email notification service implemented in (Go)
The email notification service consumes and processess messages from a queue in RabbitMQ.  The messages are json documents that include recipients list, subject string, and data used for template binding used in the body and attachments. An example of how this service would be used will be to send a thank you email after placing an order.  After an order is successfully placed, a document containing the order details is sent to the notification service. The notification service will use a handlebars template for generating the body of the email, as well as any PDF attachments.  A management interface will be available for managing data binding templates.

##RabbitMQ Processor
Clients will publish messages for sending notifications in a queue.  The serive will consume these messages and process the various notifications.

Example payload
```json
{
  "recipients": [
    "test@krillan.com",
    "test2@krillan.com"
  ],
  "subject": "Some message subject",
  "templateId": "order-thankyou",
  "data": {}
}
```

##Management API
A management API is used to list, add, edit, and delete templates.  In addition to managing templates, consumers will need to be able to view rendered templates used for editing and testing templates.

| URL | METHOD | DESCRIPTION |
| --- | --- | --- |
| /template | GET | Returns a list of templates |
| /template/:id | GET | Returns template details by id |
| /template/:id | PUT | Add or update a template |
| /template/:id | DELETE | Delete a template |
| /template/:id/view | POST | View a template using the posted data |
| /template/:id/view?format=pdf | POST | View a template using posted data as PDF attachment |
