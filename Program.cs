var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

app.MapGet("/", () => "hello_world");

app.Run();
