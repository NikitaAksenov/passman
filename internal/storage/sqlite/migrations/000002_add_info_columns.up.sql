-- Add info columns
ALTER TABLE pass ADD COLUMN created TEXT;
ALTER TABLE pass ADD COLUMN lastUpdate TEXT;
ALTER TABLE pass ADD COLUMN lastRead TEXT;

-- Set created column to current time with format "01 Jan 2025 12:34:00"
UPDATE pass
	SET created = printf("%02d %s %s %s UTC",
        CAST(strftime("%d", "now") AS INTEGER),
        CASE strftime("%m", "now")
            WHEN "01" THEN "Jan"
            WHEN "02" THEN "Feb"
            WHEN "03" THEN "Mar"
            WHEN "04" THEN "Apr"
            WHEN "05" THEN "May"
            WHEN "06" THEN "Jun"
            WHEN "07" THEN "Jul"
            WHEN "08" THEN "Aug"
            WHEN "09" THEN "Sep"
            WHEN "10" THEN "Oct"
            WHEN "11" THEN "Nov"
            WHEN "12" THEN "Dec"
        END,
        strftime("%Y", "now"),
        strftime("%H:%M:%S", "now")
);