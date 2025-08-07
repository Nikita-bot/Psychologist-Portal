document.addEventListener('DOMContentLoaded', function() {
    const editModal = new bootstrap.Modal(document.getElementById('editEmployeeModal'));
    
    document.getElementById('addEmployeeBtn').addEventListener('click', async () => {
        const employeeData = {
            name: document.getElementById('employeeName').value,
            tg: document.getElementById('employeeTelegram').value
        };
        
        if (!employeeData.name) {
            alert('Заполните ФИО сотрудника');
            return;
        }
        
        try {
            const response = await fetch('/employees', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(employeeData)
            });
            
            if (response.ok) {
                window.location.reload();
            } else {
                throw new Error('Ошибка добавления сотрудника');
            }
        } catch (error) {
            console.error('Ошибка:', error);
            alert('Не удалось добавить сотрудника');
        }
    });

    document.addEventListener('click', function(e) {
        if (e.target.closest('.btn-edit')) {
            const button = e.target.closest('.btn-edit');
            const employeeId = button.getAttribute('data-id');
            const fullName = button.getAttribute('data-fullname');
            const telegram = button.getAttribute('data-telegram');
            
            document.getElementById('editEmployeeId').value = employeeId;
            document.getElementById('editEmployeeName').value = fullName;
            document.getElementById('editEmployeeTelegram').value = telegram || '';
            
            editModal.show();
        }
    });

    document.getElementById('editEmployeeBtn').addEventListener('click', async () => {
        const employeeData = {
            id: document.getElementById('editEmployeeId').value,
            fullName: document.getElementById('editEmployeeName').value,
            telegram: document.getElementById('editEmployeeTelegram').value
        };
        
        if (!employeeData.fullName) {
            alert('Заполните ФИО сотрудника');
            return;
        }
        
        try {
            const response = await fetch(`/employees?id=${employeeData.id}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    name: employeeData.fullName,
                    tg: employeeData.telegram
                })
            });
            
            if (response.ok) {
                window.location.reload();
            } else {
                throw new Error('Ошибка сохранения изменений');
            }
        } catch (error) {
            console.error('Ошибка:', error);
            alert('Не удалось сохранить изменения');
        }
    });

    document.addEventListener('click', function(e) {
        if (e.target.closest('.btn-delete')) {
            const button = e.target.closest('.btn-delete');
            const employeeId = button.getAttribute('data-id');
            
            if (confirm('Удалить сотрудника?')) {
                fetch(`/employees?id=${employeeId}`, { method: 'DELETE' })
                    .then(response => response.ok && window.location.reload());
            }
        }
    });
});