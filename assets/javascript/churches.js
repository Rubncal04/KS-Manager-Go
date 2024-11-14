document.addEventListener("DOMContentLoaded", function () {
  const token = document.querySelector('meta[name="token"]').getAttribute("content");
  const churchContainer = document.getElementById("churchesContainer");

  const getChurches = async () => {
    try {
      const res = await fetch("/api/v1/churches", {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer " + token,
        }
      })

      const churches = await res.json();
      displayChurches(churches);
    } catch (error) {
      console.error(error);
    }
  }

  function displayChurches(churches) {
    churchContainer.innerHTML = "";
    churches.forEach((church) => {
      const churchCard = document.createElement("div");
      churchCard.classList.add("col-md-4", "mb-4");
      churchCard.innerHTML = `
            <div class="card">
                <div class="card-body">
                    <h5 class="card-title">${church.name}</h5>
                    <h6 class="card-subtitle mb-2 text-muted">${church.address}</h6>
                    <div class="mt-3">
                      <button id="show-church" data-id="${church.id}" type="button" class="btn btn-outline-primary justify-content-md-start">Show</button>
                      <button id="delete-church" data-id="${church.id}" type="button" class="btn btn-outline-danger justify-content-md-start">Delete</button>
                    </div>
                </div>
            </div>
        `;
      churchContainer.appendChild(churchCard);
    });
  }

  $("#churchesContainer").on("click", "#show-church", async function () {
    try {
      const churchId = $(this).data("id");
      const res = await fetch(`/api/v1/churches/${churchId}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          "Authorization": "Bearer " + token,
          "Accept": "text/html"
        }
      })

      window.location.href = res.url;
    } catch (error) {
      console.error(error);
    }
  });

  getChurches();
})
