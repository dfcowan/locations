namespace Locations;

public class User(long id, DateOnly startDate, DateOnly syncedThroughDate, string timeZoneId)
{
    public long Id { get; set; } = id;
    public DateOnly StartDate { get; set; } = startDate;
    public DateOnly SyncedThroughDate { get; set; } = syncedThroughDate;
    public string TimeZoneId { get; set; } = timeZoneId;
}
