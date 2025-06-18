CREATE TABLE IF NOT EXISTS logs(
        Id UInt32,
        ProjectId UInt32,
        Name String,
        Description String,
        Priority UInt32,
        Removed UInt8,
        EventTime DateTime
    ) ENGINE = MergeTree()
    ORDER BY EventTime;