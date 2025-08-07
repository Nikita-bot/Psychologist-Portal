document.addEventListener('DOMContentLoaded', function() {
    currentType = 'individual';
    
    document.querySelectorAll('.consultation-type-title').forEach(title => {
        title.addEventListener('click', function() {

            document.querySelectorAll('.consultation-type-title').forEach(t => {
                t.classList.remove('active');
            });
            
            this.classList.add('active');
            
            currentType = this.getAttribute('data-type');
            
            loadDays();
            loadSlots();
        });
    });

    function loadDays() {
        const endpoint = currentType === 'individual' 
            ? '/slots/days' 
            : '/slots_room/days';
        fetch(endpoint)
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                if (data.status === 'success') {

                    document.querySelectorAll('.weekday-btn').forEach(btn => {
                        btn.classList.remove('active');
                    });
                    console.log(data)

                    data.days.forEach(day => {
                        if (day.is_active) {
                            const btn = document.querySelector(`.weekday-btn[data-day="${day.day}"]`);
                            if (btn) btn.classList.add('active');
                        }
                    });
                }
            })
            .catch(error => {
                console.error('Error fetching days:', error);
                alert('Failed to load days: ' + error.message);
            });
    };

    function loadSlots() {
        const endpoint = currentType === 'individual' 
            ? '/slots_ind' 
            : '/slots_room';
        
        console.log(endpoint)
        const slotsList = document.getElementById('slotsList');
        slotsList.innerHTML=''
    
        fetch(endpoint)
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                if (data.status === 'success') {
                    console.log(data)
                    if (data.slots != null){
                        data.slots.forEach(slot => {
                        console.log(slot)
                        const slotCard = document.createElement('div');
                        slotCard.className = 'slot-card';
                        slotCard.innerHTML = `
                            <div class="slot-time">${slot.time}</div>
                            <div class="slot-actions">
                                <button class="btn btn-delete" data-id="${slot.id}">
                                    <i class="bi bi-trash"></i> Удалить
                                </button>
                            </div>
                        `;
                        slotsList.appendChild(slotCard);
                        });
                    }
                    else {
                        const slotCard = document.createElement('div');
                        slotCard.className = 'slot-card';
                        slotCard.innerHTML = `
                            <div class="col-12 text-center py-5">
                                <i class="bi bi-calendar-x" style="font-size: 3rem; color: #7f8c8d;"></i>
                                <h4 class="mt-3">Нет доступных слотов</h4>
                            </div>
                        `
                        slotsList.appendChild(slotCard);
                    }
                }
            })
            .catch(error => {
                console.error('Error fetching days:', error);
                alert('Failed to load days: ' + error.message);
            });
            
    };

    document.querySelectorAll('.weekday-btn').forEach(btn => {
        btn.addEventListener('click', function() {
            this.classList.toggle('active');
        });
    });

    document.getElementById('saveDaysBtn').addEventListener('click', async () => {
        const days = Array.from(document.querySelectorAll('.weekday-btn')).map(btn => ({
            day: btn.getAttribute('data-day'),
            is_active: btn.classList.contains('active')
        }));
        const endpoint = currentType === 'individual' 
            ? '/slots/days' 
            : '/slots_room/days';
        try {
            const response = await fetch(endpoint, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(days)
            });
            const result = await response.json();
            
            if (!response.ok) {
                throw new Error(result.error || 'Неизвестная ошибка сервера');
            }
            if (result.status === 'success') {
                alert('Дни успешно сохранены!');
            } else {
                throw new Error(result.message || 'Не удалось сохранить дни');
            }
        } catch (error) {
            console.error('Ошибка сохранения дней:', error);
            alert('Ошибка сохранения: ' + error.message);
        }
    });

    document.getElementById('createSlotBtn').addEventListener('click', async () => {
        const time = document.getElementById('slotTime').value;
        
        if (!time) {
            alert('Укажите время слота');
            return;
        }
        
        const endpoint = currentType === 'individual' 
            ? '/slots' 
            : '/slots_room';
        try {
            const response = await fetch(endpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    time: time
                })
            });
            
            if (response.ok) {
                loadSlots();              
            } else {
                throw new Error('Ошибка создания слота');
            }
        } catch (error) {
            console.error('Ошибка:', error);
            alert('Не удалось создать слот');
        }
    });

    document.getElementById('createSlotBtn').addEventListener('hidden.bs.modal', function () {
        document.getElementById('createSlotForm').reset();
    });

    document.addEventListener('click', function(e) {
        if (e.target.closest('.btn-delete')) {
            const button = e.target.closest('.btn-delete');
            const type = currentType;
            deleteSlot(button, type);
        }
    });

    async function deleteSlot(button, type) {
        const id = button.getAttribute('data-id');
        if (confirm('Удалить слот?')) {
            const endpoint = type === 'individual' 
                ? `/slots?id=${id}`
                : `/slots_room?id=${id}`;
            fetch(endpoint, { method: 'DELETE' })
                .then(response => response.ok && loadSlots());
        }
    }
    loadDays();
    loadSlots();
});