# Daily Tsumego (Go Problem) from goproblems.com
![image](./preview.png)

Simple widget for displaying daily tsumego (go problem) from goproblems.com website. 

```
          - type: custom-api
            title: Daily Go Problem
            cache: 6h
            url: https://www.goproblems.com/api/v2/problems/daily
            template: |
              {{ $targetCategory := "medium" }}

              {{ range .JSON.Array "entries" }}
                {{ if eq (.String "category") $targetCategory }}
                  {{ $id        := .Int "id" }}
                  {{ $imageUrl  := .String "imageUrl" }}
                  {{ $diffVal   := .Int "difficulty.value" }}
                  {{ $diffUnit  := .String "difficulty.unit" }}
                  {{ $genre     := .String "genre" }}
                  {{ $problemUrl := printf "https://www.goproblems.com/problems/%d" $id }}
                  {{ $difficulty := printf "%d %s" $diffVal $diffUnit }}

                  <div style="text-align: center;">
                    <h3 style="font-size: 1.4rem; margin: 0 0 4px 0;">
                      Medium Problem of the Day
                    </h3>
                    <div style="margin-bottom: 6px; font-size: 0.85rem; color: #888;">
                      {{ $difficulty }}    {{ $genre }}
                    </div>
                    <a href="{{ $problemUrl }}" target="_blank">
                      <img src="{{ $imageUrl }}"
                           alt="Go Problem of the Day"
                           style="max-width:100%; height:auto; border-radius: 3px; display: inline">
                    </a>
                  </div>
                {{ end }}
              {{ end }}
```
