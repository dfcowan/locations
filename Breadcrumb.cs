namespace Locations;

public class Breadcrumb(
    Guid id,
    long userId,
    double latitude,
    double longitude,
    DateTimeOffset time)
{
    public Guid Id { get; set; } = id;
    public long UserId { get; set; } = userId;
    public double Latitude { get; set; } = latitude;
    public double Longitude { get; set; } = longitude;
    public DateTimeOffset Time { get; set; } = time;
}
