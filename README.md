## WW2 veterans memorial application form

I had got tasked with making an online form platform for collecting WW2 veterans' 
data and installing a terminal for searching and reading this data near a city memorial.

https://doroga.ttnet.ru

___
### Backend

- Gin framework for API handling
- MongoDB as data storage (photos save on the local filesystem)
- ImageMagick for the images processing
- API
  - POST new application (web form handling)
  - GET (all applications with pagination or by ID)
  - GET archived applications index
  - GET archived applications for a specific date (for reporting)

___
### Frontend

Is a VueJS SPA

[sc-main.webm](https://user-images.githubusercontent.com/2155441/198938570-74a68a21-dc44-47b1-b077-f975177f2673.webm)

___
### Kiosk UI

The UI is made for a specific resolution for a kiosk mode running a browser on a street info kiosk
([Proteus Kiosk](https://porteus-kiosk.org) is used).

Has a custom virtual keyboard.

[sc-kiosk.webm](https://user-images.githubusercontent.com/2155441/198938598-37f04bd7-894b-49cb-a8c9-35711d08eb1c.webm)
