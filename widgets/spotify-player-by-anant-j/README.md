![](preview.png)

# !!!!!!! IMPORTANT - PLEASE READ
> [!CAUTION]     
> Using the play(▶) and pause(⏸) buttons sends a request directly from the browser to Spotify's Play/Pause API's respectively. These PUT requests require the bearer token to be sent in the auth header, which means the $accessToken is exposed to the client-side/browser.  
There are many security implications of this. For example, if a malicious actor get's access to this token, they can completely control and view your Spotify playlists/currently playing.  
Please use this at your own risk. [Spotify Authorization](https://developer.spotify.com/documentation/web-api/tutorials/code-flow)  
It is highly suggested to use [Glance Auth](https://github.com/glanceapp/glance/releases/tag/v0.8.0#g-rh-7). However, please note using this does not mitigate the above mentioned security risk.

Choose from one of the following otpions:  
1. Use Play/Pause in current-implementation at your own risk.
2. Move Play/Pause functionality behind an external service so accessToken is not exposed on client-side.
3. Remove Play/Pause from config below and the $accessToken will not be revealed on client-side.    

<hr>

```yaml
- type: custom-api
  hide-header: true
  cache: 1s
  frameless: true
  template: |
    {{
      $tokenRes := newRequest "https://accounts.spotify.com/api/token"
        | withHeader "Authorization" "Basic ${SPOTIFY_BTOA}"
        | withHeader "Content-Type" "application/x-www-form-urlencoded"
        | withStringBody "grant_type=refresh_token&refresh_token=${SPOTIFY_REFRESH}"
        | getResponse
    }}

    {{ if eq $tokenRes.Response.StatusCode 200 }}
      {{ $accessToken := $tokenRes.JSON.String "access_token" }}

      {{ $currentlyPlaying := newRequest "https://api.spotify.com/v1/me/player" 
          | withHeader "Authorization" (print "Bearer " $accessToken)
          | getResponse          
      }}
      {{ $isCurrentlyPlaying := $currentlyPlaying.JSON.Bool "is_playing" }}
      {{ $isDeviceActive := $currentlyPlaying.JSON.Bool "device.is_active" }}
      {{ $isPrivateSession := $currentlyPlaying.JSON.Bool "device.is_private_session" }}
      {{ $deviceName := $currentlyPlaying.JSON.String "device.name" }}
      {{ $queueRes := newRequest "https://api.spotify.com/v1/me/player/queue"
          | withHeader "Authorization" (print "Bearer " $accessToken)
          | getResponse
      }}

      {{ if eq $queueRes.Response.StatusCode 200 }}
        {{ $data := $queueRes.JSON }}
        {{ $queue := $data.Array "queue" }}

        {{ $artist := $data.String "currently_playing.artists.0.name" }}

        {{ if gt (len $artist) 0 }}
        <div class="size-h1">Spotify Now Playing</div>
        <div class="widget-content-frame flex flex-row items-center gap-20">
          <div>
            <img src="{{ $data.String "currently_playing.album.images.0.url" }}" style="border-radius: 5px; width: 5rem;" class="card">
          </div>
          <div class="flex grow flex-column justify-center pr-2">
            <div class="color-positive size-h4 font-bold">{{ $data.String "currently_playing.name" }}</div>
            <div class="size-h5">{{ $artist }}</div>
          </div>
          {{ if and $isDeviceActive (not $isPrivateSession) }}
            {{ if $isCurrentlyPlaying }}
              <div style="margin-right:10px;">
                <button style="font-size:25px;" onclick="fetch('https://api.spotify.com/v1/me/player/pause',{method:'PUT',headers:{'Authorization':'Bearer {{$accessToken}}'}}); setTimeout(function(){ location.reload(); }, 2000);">⏸</button>
                <p>{{ $deviceName }}</p>
              </div>
            {{ else }}
              <div style="margin-right:10px;">
                <button style="font-size:25px;" onclick="fetch('https://api.spotify.com/v1/me/player/play',{method:'PUT',headers:{'Authorization':'Bearer {{$accessToken}}'}}); setTimeout(function(){ location.reload(); }, 2000);">▶</button>
                <p>{{ $deviceName }}</p>
              </div>
            {{ end }}
          {{ else if $isPrivateSession }}
          <p style="margin-right:10px;">Private session is<br> active on {{ $deviceName }}</p> 
          {{ else }}
          <p style="margin-right:10px;">Device is inactive</p>
          {{ end }}
        </div>
        {{ end }}

        {{ if gt (len $queue) 0 }}
          <div class="size-h1 color-muted font-bold">Upcoming:</div>
          {{ range $i, $track := $queue }}
            {{ if lt $i 5 }}
              <div class="widget-content-frame flex flex-row items-center gap-20">
                <div>
                  <img src="{{ $track.String "album.images.0.url" }}" style="border-radius: 5px; width: 5rem;" class="card">
                </div>
                <div class="flex grow flex-column justify-center pr-2">
                  <div class="color-positive size-h4 font-bold">{{ $track.String "name" }}</div>
                  <div class="size-h5">{{ $track.String "artists.0.name" }}</div>
                </div>
              </div>
            {{ end }}
          {{ end }}
        {{ end }}

      {{ else }}
        <p>Failed to fetch queue: {{ $queueRes.Response.Status }}</p>
      {{ end }}

    {{ else }}
      <p>Failed to get token: {{ $tokenRes.Response.Status }}</p>
    {{ end }}
```

## Setup
- No external server needed.
- These tokens from Spotify are required: 
  - `client_id`
  - `client_secret`
  - `refresh`.

### ENV Variables
1. Log into [Spotify for Developers](https://developer.spotify.com/).
2. Navigate to the dashboard and create a new app, using `http://127.0.0.1:8888/callback` (or whatever) as the callback. Make sure to use https.
3. Record the `client_id` and `client_secret`. You will need this later.
4. Grab a scope variable by navigating into this link: 
`
https://accounts.spotify.com/en/authorize?client_id=<your_client_id>&response_type=code&redirect_uri=http://127.0.0.1:8888/callback&scope=user-read-currently-playing%20user-read-playback-state%20user-modify-playback-state
` 
It will redirect to a localhost url with this format: `http://127.0.0.1:8888/callback?code=<scope_variable>` - store this.  
5. [Base64 encode](https://www.base64encode.org/) the string `client_id:clientsecret`.  
6. Store this b64 encoded value as **"SPOTIFY_BTOA"** in Glance .env.  
7. [Run](https://reqbin.com/curl) the following cURL command: 
  `
  curl -H "Authorization: Basic <base64_encoded_string>"
  -d grant_type=authorization_code -d code=<scope_variable> -d redirect_uri=http://127.0.0.1:8888/callback https://accounts.spotify.com/api/token
  `
8. Record the **"refresh_token"** in the resultant `.json` file, and store as SPOTIFY_REFRESH in Glance .env.