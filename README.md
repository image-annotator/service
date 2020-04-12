# LABEL-1-BackEnd

This repository contains web server dedicated to run https://gitlab.informatika.org/if3250-labeling-project/labeling.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for public use purposes. 

### Prerequisites

What things you need to install the software and how to install them

```
latest Docker containerization tool, available for free at https://www.docker.com/
```

### Installing

A step by step series of examples that tell you how to get a development env running


0. Server can only be run in linux distro OS

1. Create this structre of folder in your home directory

   ~/go/src/gitlab.informatika.org
   
   with this command
    ```
   mkdir -p ~/go/src/gitlab.informatika.org/
    ```
2. clone this repository into 

    ~/go/src/gitlab.informatika.org
    
    your project structure should look like this now
    
     ~/go/src/gitlab.informatika.org/label-1-backend/
    
3. For first time installation run the rundatabase script with

   ```
   ./rundatabase
   ```

4. Navigate into the base directory in ~/go/src/gitlab.informatika.org/label-1-backend/base and run this command
   
   ```
   ./base migrate
   ```
   
5. Finally run the server with this command

   ```
   ./base serve
   ```
   
### Frequent Use Guide

1. After a successful installation you only need to run this command in ~/go/src/gitlab.informatika.org/label-1-backend/ to run the database
   
    ```
   ./dbstart
   ``` 
2. and navigate into base directory and run this command to run the server
   ```
   ./base serve
   ```

## Built With

* [Gin](https://github.com/gin-gonic/gin) - The web framework used
* [Docker](https://www.docker.com/) - Containerization tool

## Authors

* **Rayza Mahendra** - *Backend Programmer* - [rayzamgh](https://github.com/rayzamgh)
* **Edward Alexander** - *Frontend Programmer* - [rayzamgh](https://github.com/rayzamgh)
* **Ahmad Rizal Alifio** - *Backend Programmer* - [ARAlifio](https://github.com/ARAlifio)
* **Nurdin** - *Frontend Programmer* - [rayzamgh](https://github.com/rayzamgh)
* **Eka Sunandika** - *Frontend Programmer* - [rayzamgh](https://github.com/rayzamgh)

## License

This project is licensed under the MIT License 

MIT License

Copyright (c) [2020] [Label]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.


## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc
