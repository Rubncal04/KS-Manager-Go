const form = document.getElementById('login-form');

form.addEventListener('submit', async function (event) {
  event.preventDefault();

  const formData = new FormData(form);

  const data = {};
  formData.forEach((value, key) => {
    data[key] = value;
  });

  try {
    const response = await fetch('/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    if (response.ok) {
      const result = await response.json();
      console.log('Login successful:', result);
    } else {
      console.error('Login failed:', response.statusText);
    }
  } catch (error) {
    console.error('Error during login:', error);
  }
});
