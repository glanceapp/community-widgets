This widget tracks upcoming holidays, dynamic family birthdays, and trip departures while automatically adjusting text colors as event dates approach. Curently, it tracks Christmas, New Years Day, a vacation date, and 2 birthdays of your choice.

## Color Meanings
## green is either green, or your standard theme accent color.
---

* **Conditional Status Indicators:** The colors change as you get closer to the date:
    * **🔴 Critical Alert (< 7 days):** 
    * **🟡 Warning Window (< 30 days):**
    * **🟢 General State (> 30 days):**
---

## Installation & Configuration

Copy this into a file of your choice, I named it "countdowns.yml" or you can put it straight in your pages. I call mine using an $include: countdowns.yml in my pages that i want it on.:

```yaml
- type: custom-api
  title: Countdowns
  cache: 1h
  template: |
    <div style="display:grid; grid-template-columns:1fr 1fr; gap:12px;">

      {{ \$christmas := parseTime "2006-01-02" "{CHRISTMAS}" }}
      {{ christmasDays := div (christmas.Sub now).Hours 24.0 }}
      <div style="text-align:center;">
        <div class="size-h5">🎄 Christmas</div>
        {{ if le \$christmasDays 7.0 }}
          <div class="size-h3" style="color:var(--color-negative);">{{ printf "%.0f" \$christmasDays }} days</div>
        {{ else if le \$christmasDays 30.0 }}
          <div class="size-h3" style="color:var(--color-primary);">{{ printf "%.0f" \$christmasDays }} days</div>
        {{ else }}
          <div class="size-h3 color-highlight">{{ printf "%.0f" \$christmasDays }} days</div>
        {{ end }}
      </div>

      {{ \$newyear := parseTime "2006-01-02" "{NEW_YEARS_DAY)" }}
      {{ newyearDays := div (newyear.Sub now).Hours 24.0 }}
      <div style="text-align:center;">
        <div class="size-h5">🎆 New Year</div>
        {{ if le \$newyearDays 7.0 }}
          <div class="size-h3" style="color:var(--color-negative);">{{ printf "%.0f" \$newyearDays }} days</div>
        {{ else if le \$newyearDays 30.0 }}
          <div class="size-h3" style="color:var(--color-primary);">{{ printf "%.0f" \$newyearDays }} days</div>
        {{ else }}
          <div class="size-h3 color-highlight">{{ printf "%.0f" \$newyearDays }} days</div>
        {{ end }}
      </div>

      {{ vacation := parseTime "2006-01-02" "{VACATION}" }}
      {{ vacationDays := div (vacation.Sub now).Hours 24.0 }}
      <div style="text-align:center;">
        <div class="size-h5">✈️ Next Vacation</div>
        {{ if le \$vacationDays 7.0 }}
          <div class="size-h3" style="color:var(--color-negative);">{{ printf "%.0f" \$vacationDays }} days</div>
        {{ else if le \$vacationDays 30.0 }}
          <div class="size-h3" style="color:var(--color-primary);">{{ printf "%.0f" \$vacationDays }} days</div>
        {{ else }}
          <div class="size-h3 color-highlight">{{ printf "%.0f" \$vacationDays }} days</div>
        {{ end }}
      </div>

      {{ birthday1 := parseTime "2006-01-02" "{BIRTHDAY_1}" }}
      {{ birthday1Days := div (birthday1.Sub now).Hours 24.0 }}
      <div style="text-align:center;">
        <div class="size-h5">🎂 First person's Birthday</div>
        {{ if le \$birthday1Days 7.0 }}
          <div class="size-h3" style="color:var(--color-negative);">{{ printf "%.0f" \$birthday1Days }} days</div>
        {{ else if le \$birthday1Days 30.0 }}
          <div class="size-h3" style="color:var(--color-primary);">{{ printf "%.0f" \$birthday1Days }} days</div>
        {{ else }}
          <div class="size-h3 color-highlight">{{ printf "%.0f" \$birthday1Days }} days</div>
        {{ end }}
      </div>

      {{ birthday2 := parseTime "2006-01-02" "{BIRTHDAY_2}" }}
      {{ birthday2Days := div (birthday2.Sub now).Hours 24.0 }}
      <div style="text-align:center;">
        <div class="size-h5">🎂 Second person's Birthday</div>
        {{ if le \$birthday2Days 7.0 }}
          <div class="size-h3" style="color:var(--color-negative);">{{ printf "%.0f" \$birthday2Days }} days</div>
        {{ else if le \$birthday2Days 30.0 }}
          <div class="size-h3" style="color:var(--color-primary);">{{ printf "%.0f" \$birthday2Days }} days</div>
        {{ else }}
          <div class="size-h3 color-highlight">{{ printf "%.0f" \$birthday2Days }} days</div>
        {{ end }}
      </div>
    </div>
```

---

## 🛠️ Required Environment Variables

To render any of the fields properly, ensure the variables are using this standard **`YYYY-MM-DD`** layout. Anything else will result in dates being wrong on the countdown.


### Standard Environment File (`.env`)
### These environment variables are needed for all the dates. Note that you will have to change the dates per year.
```env
CHRISTMAS=2026-12-25
NEW_YEARS_DAY=2026-01-01
VACATION=2026-08-27
BIRTHDAY_1=2026-11-20
BIRTHDAY_2=2027-03-05
```

---

