document.addEventListener('DOMContentLoaded', () => {
    // Show/hide sections
    document.getElementById('showRegister').addEventListener('click', (e) => {
        e.preventDefault();
        document.getElementById('authSection').style.display = 'none';
        document.getElementById('registerSection').style.display = 'block';
        document.getElementById('calcSection').style.display = 'none';
    });

    document.getElementById('showLogin').addEventListener('click', (e) => {
        e.preventDefault();
        document.getElementById('authSection').style.display = 'block';
        document.getElementById('registerSection').style.display = 'none';
        document.getElementById('calcSection').style.display = 'none';
    });

    // Handle registration
    document.getElementById('registerForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const login = document.getElementById('registerLoginInput').value;
        const password = document.getElementById('registerPasswordInput').value;

        try {
            const response = await fetch('/api/v1/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    login: login,
                    password: password
                })
            });

            if (response.ok) {
                alert('Регистрация успешна! Теперь вы можете войти.');
                document.getElementById('showLogin').click();
            } else {
                const data = await response.json();
                alert(data.error || 'Ошибка при регистрации');
            }
        } catch (error) {
            console.error('Error:', error);
            alert('Ошибка при регистрации');
        }
    });

    // Handle login
    document.getElementById('loginForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const login = document.getElementById('loginInput').value;
        const password = document.getElementById('passwordInput').value;

        try {
            const response = await fetch('/api/v1/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    login: login,
                    password: password
                })
            });

            if (response.ok) {
                const data = await response.json();
                localStorage.setItem('token', data.token);
                document.getElementById('authSection').style.display = 'none';
                document.getElementById('registerSection').style.display = 'none';
                document.getElementById('calcSection').style.display = 'block';
                loadExpressions();
            } else {
                const data = await response.json();
                alert(data.error || 'Ошибка при входе');
            }
        } catch (error) {
            console.error('Error:', error);
            alert('Ошибка при входе');
        }
    });

    // Handle expression submission
    document.getElementById('expressionForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const input = document.getElementById('expressionInput');
        const token = localStorage.getItem('token');

        if (!token) {
            alert('Необходимо войти в систему');
            return;
        }

        try {
            const response = await fetch('/api/v1/calculate', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    expression: input.value
                })
            });

            if (response.ok) {
                input.value = '';
                loadExpressions();
            } else {
                const data = await response.json();
                alert(data.error || 'Ошибка при отправке выражения');
            }
        } catch (error) {
            console.error('Error:', error);
            alert('Ошибка при отправке выражения');
        }
    });

    async function loadExpressions() {
        const list = document.getElementById('expressionsList');
        const token = localStorage.getItem('token');

        if (!token) {
            list.innerHTML = '<div class="error">Необходимо войти в систему</div>';
            return;
        }

        list.innerHTML = '<div class="loading">Загрузка...</div>';

        try {
            const response = await fetch('/api/v1/expressions', {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });
            const data = await response.json();

            if (response.ok) {
                list.innerHTML = data.expressions.map(expr => `
                    <div class="expression-item">
                        <div>
                            <strong>${expr.expression}</strong>
                            <div class="status ${expr.status}">${expr.status}</div>
                        </div>
                        <div>${expr.result || ''}</div>
                    </div>
                `).join('');
            } else {
                list.innerHTML = '<div class="error">Ошибка загрузки данных</div>';
            }
        } catch (error) {
            list.innerHTML = '<div class="error">Ошибка загрузки данных</div>';
        }
    }

    // Check if user is already logged in
    const token = localStorage.getItem('token');
    if (token) {
        document.getElementById('authSection').style.display = 'none';
        document.getElementById('registerSection').style.display = 'none';
        document.getElementById('calcSection').style.display = 'block';
        loadExpressions();
    }
});