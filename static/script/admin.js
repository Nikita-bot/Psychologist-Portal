async function cancelEvent(button) {
    const eventId = button.getAttribute('data-event-id');
    if (!eventId) {
        console.error('Event ID not found');
        return;
    }

    // Показываем подтверждение
    if (!confirm('Вы уверены, что хотите отменить эту консультацию?')) {
        return;
    }

    try {
        const response = await fetch(`/admin/event?id=${eventId}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Ошибка при отмене консультации');
        }

        // Если успешно - обновляем страницу
        window.location.reload();
    } catch (error) {
        console.error('Error:', error);
        alert(error.message);
    }
}