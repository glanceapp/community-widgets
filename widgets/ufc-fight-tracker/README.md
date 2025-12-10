# Preview
![UFC Fight Tracker Widget Preview](preview.png)

# Notes
- In the preview this widget was used in a column between two small side columns. All other configurations have not been tested.
- No additional configuration is necessary for this widget.

# Update Frequency
Every 10 minutes by default (can be changed).

# YAML
```yaml
- type: custom-api
          title: UFC Fight Tracker
          cache: 10m
          url: https://site.api.espn.com/apis/site/v2/sports/mma/ufc/scoreboard
          template: |
            <ul class="list collapsible-container flex flex-column gap-10" data-collapse-after="3" style="list-style: none; padding: 0; margin: 0;">
              
              {{ $events := .JSON.Array "events" }}
              
              {{ if $events }}
                {{ range $events }}
                  {{ $status := .String "status.type.state" }}
                  {{ $dateStr := .String "date" }}
                  {{ $eventName := .String "name" }}
                  
                  {{ range .Array "competitions" }}
                    
                    {{/* VENUE LOGIC: Combine Name + City/State/Country */}}
                    {{ $vName := .String "venue.fullName" }}
                    {{ $vCity := .String "venue.address.city" }}
                    {{ $vCountry := .String "venue.address.country" }}
                    
                    {{/* Build location string based on what is available */}}
                    {{ $location := $vName }}
                    {{ if $vCity }}
                      {{ $location = print $location " • " $vCity }}
                      {{ if $vCountry }}
                         {{ $location = print $location ", " $vCountry }}
                      {{ end }}
                    {{ end }}

                    <li>
                      <div style="background: rgba(255,255,255,0.03); border-radius: 8px; padding: 12px; border: 1px solid rgba(255,255,255,0.05);">
                        
                        <div class="flex justify-between items-center" style="margin-bottom: 12px; border-bottom: 1px solid rgba(255,255,255,0.1); padding-bottom: 8px;">
                          <div class="text-truncate" style="font-weight: 700; color: var(--color-highlight); font-size: 0.9em; max-width: 60%;">
                            {{ $eventName }}
                          </div>
                          <div style="font-size: 0.75em; font-family: monospace; text-align: right;">
                            {{ if eq $status "in" }}
                              <span style="color: #ef4444; font-weight: bold; animation: pulse 2s infinite;">🔴 LIVE R{{ .String "status.period" }}</span>
                            {{ else if eq $status "post" }}
                              <span style="color: var(--color-subdue);">FINAL</span>
                            {{ else }}
                              {{ if $dateStr }}
                                {{ $t := parseTime "2006-01-02T15:04Z" $dateStr }}
                                📅 {{ $t.Local.Format "02. Jan • 15:04" }}
                              {{ end }}
                            {{ end }}
                          </div>
                        </div>

                        <div style="display: grid; grid-template-columns: 1fr 40px 1fr; align-items: center;">
                          {{ range $index, $competitor := .Array "competitors" }}
                            
                            {{ if eq $index 1 }}
                              <div style="text-align: center; font-size: 0.75em; color: var(--color-subdue); font-weight: bold; grid-column: 2;">VS</div>
                            {{ end }}

                            <div style="grid-column: {{ if eq $index 0 }}1{{ else }}3{{ end }}; display: flex; align-items: center; {{ if eq $index 0 }}justify-content: flex-end; text-align: right;{{ else }}justify-content: flex-start; text-align: left;{{ end }}">
                              
                              <div style="display: flex; align-items: center; gap: 12px; {{ if eq $index 1 }}flex-direction: row-reverse;{{ end }}">
                                <div style="min-width: 0;">
                                  <div style="display: flex; align-items: center; gap: 6px; {{ if eq $index 0 }}justify-content: flex-end;{{ else }}justify-content: flex-start;{{ end }}">
                                    {{ if and (eq $index 0) (.String "athlete.flag.href") }}
                                       <img src="{{ .String "athlete.flag.href" }}" style="width: 14px; height: 10px; opacity: 0.8; border-radius: 2px;">
                                    {{ end }}
                                    
                                    <span class="text-truncate" style="font-weight: bold; font-size: 0.9em; {{ if .Bool "winner" }}color: #22c55e;{{ end }}">
                                      {{ .String "athlete.displayName" }}
                                    </span>

                                    {{ if and (eq $index 1) (.String "athlete.flag.href") }}
                                       <img src="{{ .String "athlete.flag.href" }}" style="width: 14px; height: 10px; opacity: 0.8; border-radius: 2px;">
                                    {{ end }}
                                  </div>
                                  
                                  <div style="font-size: 0.75em; opacity: 0.6; margin-top: 2px;">
                                    {{ if .String "curatedRank.current" }}<span style="color:var(--color-highlight)">#{{ .String "curatedRank.current" }}</span> • {{ end }}
                                    {{ .String "records.0.summary" }}
                                  </div>
                                </div>

                                <div style="position: relative; width: 42px; height: 42px; shrink: 0;">
                                  <img src="{{ .String "athlete.headshot.href" }}" 
                                       style="width: 100%; height: 100%; object-fit: cover; border-radius: 50%; background: rgba(0,0,0,0.3); border: 1px solid rgba(255,255,255,0.1);"
                                       onerror="this.style.opacity='0'">
                                </div>
                              </div>
                            </div>
                          {{ end }}
                        </div>
                        
                        <div style="margin-top: 10px; font-size: 0.7em; text-align: center; opacity: 0.4;">
                          📍 {{ $location }}
                        </div>

                      </div>
                    </li>
                  {{ end }}
                {{ end }}
              {{ else }}
                <div class="color-subdue">No scheduled fights found.</div>
              {{ end }}
            </ul>
```
