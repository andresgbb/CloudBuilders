
const input = document.getElementById("input");
const output = document.getElementById("output");


input.addEventListener("keydown", async (e) => {
  if (e.key === "Enter") {
    const command = input.value.trim();

    if (command === "") return; 

   
    output.innerHTML += `$ ${command}\n`;
    input.value = "";


    if (command === "clear") {
      output.innerHTML = "";
      return;
    }

   
    const res = await fetch("/api/command", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ command }),
    });

    
    const data = await res.json();

    
    output.innerHTML += data.result + "\n";

    
    window.scrollTo(0, document.body.scrollHeight);
  }
});
