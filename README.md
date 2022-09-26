# GophKeeper

## Available methods to client
<!-- TOC -->
* [GophKeeper](#gophkeeper)
  * [Available methods to client](#available-methods-to-client)
    * [Login](#login)
    * [Register](#register)
    * [Logout](#logout)
    * [Delete logged user](#delete-logged-user)
    * [Get list of secret type](#get-list-of-secret-type)
    * [Store Login/Pass](#store-loginpass)
    * [Store Text](#store-text)
    * [Store Card](#store-card)
    * [Get secret](#get-secret)
    * [Store binary secret](#store-binary-secret)
    * [Get binary secret](#get-binary-secret)
    * [Delete secret](#delete-secret)
    * [Edit secret](#edit-secret)
    * [Get list of secret by provided type](#get-list-of-secret-by-provided-type)
    * [Exit](#exit)
<!-- TOC -->


### Login

`login %username% %password%`

### Register

`register %username% %password%`

### Logout

`logout`

### Delete logged user

`delete-user`

### Get list of secret type

`types`

### Store Login/Pass

`create-auth %title% %login% %pass%`

### Store Text

`create-text %title% %text%`

### Store Card

`create-card %title% %cardNumber% %cvv% %dueDate%`

### Get secret

`get-secret %id%`

### Store binary secret

`create-binary %title% %absolutePath%`

### Get binary secret

`get-secret-binary %id% %absolutePath%`

### Delete secret

`delete-secret id`

### Edit secret

`edit-secret %id% %title% %typeId% %fields...%`

> Number of fields needed to be passed in, is based on typeId of record.

> If local storage copy of data is not in sync with server, on attempting to update data on the server error
> will be thrown and re-sync will be started. To always rewrite data on the server ignore in sync local storage or not
> you can pass -f or --force flag.

### Get list of secret by provided type

`get-secrets-by-type %typeId%`

### Exit

`exit`