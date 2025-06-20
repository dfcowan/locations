using System.Text.Json;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;

namespace Locations.Controllers;

[Route("api/traccar")]
[ApiController]
public class TraccarController(LocationsContext context) : ControllerBase
{
    private readonly LocationsContext _context = context;

    [HttpPost]
    public async Task<IActionResult> PostBreadcrumbAsync([FromBody] JsonElement jsonElement)
    {
        OsmAndRequest? osmAndRequest = null;
        try
        {
            osmAndRequest = JsonSerializer.Deserialize<OsmAndRequest>(jsonElement);
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Exception - {ex.Message} at {ex.StackTrace} while parsing {JsonSerializer.Serialize(jsonElement)}");
            return BadRequest("Couldn't deserialize input JSON");
        }

        if (osmAndRequest == null)
        {
            return BadRequest("Couldn't deserialize input JSON. Null...");
        }

        long id = long.Parse(osmAndRequest.DeviceId);
        double lat = osmAndRequest.Location.Coords.Latitude;
        double lon = osmAndRequest.Location.Coords.Longitude;
        lat = Math.Round(lat, 5);
        lon = Math.Round(lon, 5);
        DateTimeOffset now = DateTimeOffset.UtcNow;
        var time = new DateTimeOffset(osmAndRequest.Location.Timestamp, TimeSpan.Zero);
        double accuracy = osmAndRequest.Location.Coords.Accuracy;

        if (accuracy < 0 || accuracy > 77.55)
        {
            Console.WriteLine($"accuracy is invalid - {accuracy} {id} {lat} {lon} {time}");
            Console.WriteLine(JsonSerializer.Serialize(jsonElement));
            return BadRequest($"accuracy is invalid - {accuracy}");
        }

        if (time > now.AddHours(1))
        {
            Console.WriteLine($"time is too far in the future - {time}");
            return BadRequest($"time is too far in the future - {time}");
        }

        long userId = id;
        User? user = await _context.Users.Where(u => u.Id == userId).FirstOrDefaultAsync();
        if (user == null)
        {
            Console.WriteLine($"user {userId} not found");
            return NotFound($"user {userId} not found");
        }

        Breadcrumb breadcrumb = new(
            id: Guid.NewGuid(),
            userId: userId,
            latitude: lat,
            longitude: lon,
            time: time);
        _context.Breadcrumbs.Add(breadcrumb);

        TimeZoneInfo tz = TimeZoneInfo.FindSystemTimeZoneById(user.TimeZoneId);
        DateTimeOffset usersTime = TimeZoneInfo.ConvertTime(time, tz);
        DateOnly userDate = DateOnly.FromDateTime(usersTime.Date);
        if (userDate > user.SyncedThroughDate)
        {
            user.SyncedThroughDate = userDate;
        }

        await _context.SaveChangesAsync();

        return Created();
    }
}