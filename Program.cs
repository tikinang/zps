var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

var x = Environment.GetEnvironmentVariable("MY_X_VAR");
app.MapGet("/", () => "hello_world" + x);

app.Run();
