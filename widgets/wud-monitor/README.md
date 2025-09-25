# What's Up Docker Monitor - MENU
The **What's Up Docker Monitor** widget uses the WUD API to display your Docker containers and check if updates are available. It shows the status of each container and allows you to toggle between viewing all containers or only those that need an update. It allows you to choose between forcing wud to refresh container data or fetch existing one. 

---

## Available Versions
I decided to keep old version of wud, because I know that for some people it might look bad or be overcomplicated. I'll try to maintain this version, but I don't promise anything. If you want to request some changes you can find contact info on the bottom. 

### 2.1 - [**Visual Improvment and Overall Info**](wud-main/README.md)
- **Update Description**: Fixed displaying sha.
- **Features**:
  - Better look achived by using CSS
  - Overall info
  - A lot of bug fixes

    [![2.0 Preview](wud-main/wud_main_preview.png)](wud-main/wud_main_preview.png)

### 1.1 - [**Old Release**](wud-old/README.md)
- **Description**: In v1.1, the widget now uses a `GET` request instead of `POST` to fetch container data, improving performance. 
- **Features**:
  - Simple interface & structure
  - Displays Docker container statuses and update information
    
    [![1.1 Preview](wud-old/preview_2.png)](wud-old/preview_2.png)
 
## Environment Variables

`WUD_URL` – The URL of the What's Up Docker server

Template: `WUD_URL=ip:port` – You can also just replace the code variable for it to work.

To grab containers regardless of their state, I also recommend adding this to your WUD environment:

```txt

- WUD_WATCHER_LOCAL_WATCHALL=true

```
Please remember to restart your services after applying environment variables.

---

**Made by:** Artur Flis  
**Contact:** @blue.dev on the project’s Discord
