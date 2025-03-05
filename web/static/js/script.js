document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('expressionForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        const input = document.getElementById('expressionInput');

        try {
            const response = await fetch('/api/v1/calculate', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    expression: input.value
                })
            });

            if (response.ok) {
                input.value = '';
                loadExpressions();
            } else {
                alert('Ошибка при отправке выражения');
            }
        } catch (error) {
            console.error('Error:', error);
        }
    });

    async function loadExpressions() {
        const list = document.getElementById('expressionsList');
        list.innerHTML = '<div class="loading">Загрузка...</div>';

        try {
            const response = await fetch('/api/v1/expressions');
            const data = await response.json();

            list.innerHTML = data.expressions.map(expr => `
                <div class="expression-item">
                    <div>
                        <strong>${expr.expression}</strong>
                        <div class="status ${expr.status}">${expr.status}</div>
                    </div>
                    <div>${expr.result || ''}</div>
                </div>
            `).join('');
        } catch (error) {
            list.innerHTML = '<div class="error">Ошибка загрузки данных</div>';
        }
    }

    loadExpressions();
});