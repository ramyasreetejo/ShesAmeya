document.addEventListener("DOMContentLoaded", () => {
    const messageInput = document.getElementById("message");
    const responseBox = document.getElementById("response");
    const sendBtn = document.getElementById("sendBtn");
  
    sendBtn.addEventListener("click", async () => {
      const message = messageInput.value.trim();
      if (!message) return;
  
      responseBox.innerText = "Thinking...";
  
      try {
        const ipRes = await fetch("https://api.ipify.org?format=json");
        const ipData = await ipRes.json();
  
        const res = await fetch("/api/chat", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ message, ip: ipData.ip }),
        });
  
        const data = await res.json();
        responseBox.innerHTML = data.reply
          ? marked.parse(data.reply)
          : "No response.";
      } catch (e) {
        responseBox.innerText = "Error: " + e.message;
      }
    });
  });
  