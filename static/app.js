// Obtenemos referencias a los elementos HTML
const input = document.getElementById("input");
const output = document.getElementById("output");

// Escuchamos cuando el usuario presiona Enter
input.addEventListener("keydown", async (e) => {
  if (e.key === "Enter") {
    const command = input.value.trim();

    if (command === "") return; // Evitar comandos vacíos

    // Mostramos el comando en pantalla
    output.innerHTML += `$ ${command}\n`;
    input.value = "";

    // Si el comando fue "clear", limpiamos y salimos
    if (command === "clear") {
      output.innerHTML = "";
      return;
    }

    // Enviamos el comando al backend Go
    const res = await fetch("/api/command", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ command }),
    });

    // Esperamos la respuesta JSON
    const data = await res.json();

    // Mostramos el resultado en la terminal
    output.innerHTML += data.result + "\n";

    // Hacemos scroll hacia abajo automáticamente
    window.scrollTo(0, document.body.scrollHeight);
  }
});
