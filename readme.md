# GophKeeper

## Available methods to client

### Login

`login username password`

### Register

`register username password`

### Logout

`logout`

### Delete logged user

`delete-user`

### Get list of secret type

`secret-types typeId`

### Store Login/Pass

`create-auth-secret title login pass`

### Store Text

`create-text-secret title text`

### Store Card

`create-card-secret title cardNumber cvv dueDate`

### Get secret

`get-secret id`

### Store binary secret

`create-binary-secret title absolutePath`

### Get binary secret

`get-binary-secret id absolutePath`

### Delete secret

`delete-secret id`

### Edit secret

`edit-secret id title typeId fields...`
### Get list of secret by provided type

`get-secret-list-by-type typeId`

### Exit

`exit`