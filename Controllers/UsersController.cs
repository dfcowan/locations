using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;

namespace Locations.Controllers;

[Route("api/Users")]
[ApiController]
public class UsersController(LocationsContext context) : ControllerBase
{
    private readonly LocationsContext _context = context;

    [HttpGet("{userId}/counts")]
    public async Task<IActionResult> GetCountsAsync(
        [FromRoute] long userId,
        [FromQuery] DateTime? startDate,
        [FromQuery] DateTime? endDate)
    {
        if (startDate == null)
        {
            return BadRequest("startDate is required");
        }
        if (endDate == null)
        {
            return BadRequest("endDate is required");
        }

        User? user = await _context.Users.Where(u => u.Id == userId).FirstOrDefaultAsync();
        if (user == null)
        {
            return NotFound($"user {userId} not found");
        }

        TimeZoneInfo tz = TimeZoneInfo.FindSystemTimeZoneById(user.TimeZoneId);

        DateTimeOffset startTime = new(startDate.Value, tz.GetUtcOffset(startDate.Value));
        startTime = startTime.ToUniversalTime();
        DateTimeOffset endTime = new(endDate.Value, tz.GetUtcOffset(endDate.Value));
        endTime = endTime.ToUniversalTime();

        var counts = await _context.Breadcrumbs
            .Where(bc => bc.UserId == userId && bc.Time >= startTime && bc.Time <= endTime)
            .GroupBy(bc => new { bc.Latitude, bc.Longitude })
            .Select(g => new { latitude = g.Key.Latitude, longitude = g.Key.Longitude, count = g.Count() })
            .ToListAsync();

        return Ok(counts);
    }

    [HttpGet("{userId}/sync")]
    public async Task<IActionResult> GetSyncAsync([FromRoute] long userId)
    {
        User? user = await _context.Users.Where(u => u.Id == userId).FirstOrDefaultAsync();
        if (user == null)
        {
            return NotFound($"user {userId} not found");
        }

        return Ok(user);
    }
}