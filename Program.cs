var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

app.MapGet("/", () => "hello_world" + x);
app.MapGet("/version", () => Environment.Version.ToString());
app.MapGet("/var", () => Environment.GetEnvironmentVariable("MY_X_VAR"));

app.Run();
