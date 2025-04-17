using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;

namespace Locations.Controllers;

[Route("api/traccar")]
[ApiController]
public class TraccarController(LocationsContext context) : ControllerBase
{
    private readonly LocationsContext _context = context;

    [HttpPost]
    public async Task<IActionResult> PostBreadcrumbAsync(
        [FromQuery] long? id,
        [FromQuery] double? lat,
        [FromQuery] double? lon,
        [FromQuery] long? timestamp,
        [FromQuery] long? accuracy
        )
    {
        if (id == null)
        {
            return BadRequest("id is required");
        }
        if (lat == null)
        {
            return BadRequest("lat is required");
        }
        if (lon == null)
        {
            return BadRequest("lon is required");
        }
        if (timestamp == null)
        {
            return BadRequest("timestamp is required");
        }
        if (accuracy == null)
        {
            return BadRequest("accuracy is required");
        }

        if (accuracy < 0 || accuracy > 75)
        {
            return BadRequest($"accuracy is invalid - {accuracy}");
        }

        DateTimeOffset now = DateTimeOffset.UtcNow;
        Console.WriteLine(now.ToUnixTimeSeconds());
        DateTimeOffset time = DateTimeOffset.FromUnixTimeSeconds(timestamp.Value);

        if (time > now.AddHours(1))
        {
            return BadRequest($"time is too far in the future - {time}");
        }

        long userId = id.Value;
        User? user = await _context.Users.Where(u => u.Id == userId).FirstOrDefaultAsync();
        if (user == null)
        {
            return NotFound($"user {userId} not found");
        }

        Breadcrumb breadcrumb = new(
            id: Guid.NewGuid(),
            userId: userId,
            latitude: lat.Value,
            longitude: lon.Value,
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