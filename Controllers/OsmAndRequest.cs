using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

public class OsmAndRequest(
    OsmAndRequestLocation location,
    string deviceId) // Primary constructor with parameters
{
    [Required] // Mark as required
    [JsonPropertyName("location")]
    public OsmAndRequestLocation Location { get; set; } = location;

    [Required] // Mark as required
    [JsonPropertyName("device_id")]
    public string DeviceId { get; set; } = deviceId;
}

public class OsmAndRequestLocation(
    OsmAndRequestCoords coords,
    Dictionary<string, object> extras,
    bool isMoving,
    string @event, // @event to avoid keyword conflict
    double odometer,
    string rawDataString,
    OsmAndRequestActivity activity,
    OsmAndRequestBattery battery,
    DateTime timestamp) // Primary constructor with parameters
{
    [Required] // Mark as required
    [JsonPropertyName("coords")]
    public OsmAndRequestCoords Coords { get; set; } = coords;

    [Required] // Mark as required
    [JsonPropertyName("extras")]
    public Dictionary<string, object> Extras { get; set; } = extras; // Use Dictionary<string, object> for flexible key-value pairs

    [Required] // Mark as required
    [JsonPropertyName("is_moving")]
    public bool IsMoving { get; set; } = isMoving;

    [Required] // Mark as required
    [JsonPropertyName("event")]
    public string Event { get; set; } = @event;

    [Required] // Mark as required
    [JsonPropertyName("odometer")]
    public double Odometer { get; set; } = odometer;

    [Required] // Mark as required
    [JsonPropertyName("_")] // Special property name, use JsonPropertyName
    public string RawDataString { get; set; } = rawDataString;

    [Required] // Mark as required
    [JsonPropertyName("activity")]
    public OsmAndRequestActivity Activity { get; set; } = activity;

    [Required] // Mark as required
    [JsonPropertyName("battery")]
    public OsmAndRequestBattery Battery { get; set; } = battery;

    [Required] // Mark as required
    [JsonPropertyName("timestamp")]
    public DateTime Timestamp { get; set; } = timestamp;
}

public class OsmAndRequestCoords(
    double heading,
    double speed,
    double latitude,
    double longitude,
    double accuracy,
    double altitude) // Primary constructor with parameters
{
    [Required] // Mark as required
    [JsonPropertyName("heading")]
    public double Heading { get; set; } = heading;

    [Required] // Mark as required
    [JsonPropertyName("speed")]
    public double Speed { get; set; } = speed;

    [Required] // Mark as required
    [JsonPropertyName("latitude")]
    public double Latitude { get; set; } = latitude;

    [Required] // Mark as required
    [JsonPropertyName("longitude")]
    public double Longitude { get; set; } = longitude;

    [Required] // Mark as required
    [JsonPropertyName("accuracy")]
    public double Accuracy { get; set; } = accuracy;

    [Required] // Mark as required
    [JsonPropertyName("altitude")]
    public double Altitude { get; set; } = altitude;
}

public class OsmAndRequestActivity(string type) // Primary constructor with parameters
{
    [Required] // Mark as required
    [JsonPropertyName("type")]
    public string Type { get; set; } = type;
}

public class OsmAndRequestBattery(
    double level,
    bool isCharging) // Primary constructor with parameters
{
    [Required] // Mark as required
    [JsonPropertyName("level")]
    public double Level { get; set; } = level;

    [Required] // Mark as required
    [JsonPropertyName("is_charging")]
    public bool IsCharging { get; set; } = isCharging;
}