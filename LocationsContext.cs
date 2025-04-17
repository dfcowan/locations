using Microsoft.EntityFrameworkCore;

namespace Locations;

public class LocationsContext(DbContextOptions<LocationsContext> options) : DbContext(options)
{
    public DbSet<User> Users { get; set; }
    public DbSet<Breadcrumb> Breadcrumbs { get; set; }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<Breadcrumb>()
            .HasIndex(b => new { b.UserId, b.Time });
    }
}
