document.addEventListener('DOMContentLoaded', function() {
    
    let currentStep = 1;
    const totalSteps = 7;

    fetch('/slots/days')
        .then(response => {
            if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
            return response.json();
        })
        .then(data => {
            if (data.status !== 'success') throw new Error('Invalid response format');
            console.log(data)
            // Создаем массив активных дней недели (0-воскресенье, 1-понедельник и т.д.)
            const activeDays = data.days
                .filter(day => day.is_active)
                .map(day => {
                    const daysMap = {
                        'sunday': 0,
                        'monday': 1,
                        'tuesday': 2,
                        'wednesday': 3,
                        'thursday': 4,
                        'friday': 5,
                        'saturday': 6
                    };
                    return daysMap[day.day.toLowerCase()];
                });
            console.log(activeDays)
            
            const datePicker = flatpickr("#date", {
                minDate: "today",
                maxDate: new Date().fp_incr(90), // +3 месяца
                disable: [
                    function(date) {
                        return !activeDays.includes(date.getDay());
                    }
                ],
                locale: "ru",
                dateFormat: "Y-m-d",
                onChange: function(selectedDates) {
                    console.log("Выбрана дата:", selectedDates[0]);
                },
                onDayCreate: function(dObj, dStr, fp, dayElem) {
                    // Помечаем рабочие дни
                    if (activeDays.includes(dayElem.dateObj.getDay())) {
                        dayElem.classList.add("work-day");
                    }
                }
            });
            
        })
        .catch(error => {
            console.error('Error loading active days:', error);
            alert('Не удалось загрузить доступные дни. Пожалуйста, попробуйте позже.');
            
            // Фолбэк - базовый календарь если API не отвечает
            flatpickr("#date", {
                minDate: "today",
                maxDate: new Date().fp_incr(90),
                locale: "ru",
                dateFormat: "Y-m-d"
            });
        });

    const form = document.getElementById('step-form');
    const progressBar = document.querySelector('.progress-bar');
    const prevButtons = document.querySelectorAll('.btn-prev');
    const nextButtons = document.querySelectorAll('.btn-next');
    const steps = document.querySelectorAll('.step');
    const formSteps = document.querySelectorAll('.form-step');
    const restartBtn = document.getElementById('restart-btn');
    
    const radioOptions = document.querySelectorAll('.radio-option');
    radioOptions.forEach(option => {
        option.addEventListener('click', function() {
            radioOptions.forEach(opt => opt.classList.remove('selected'));
            
            this.classList.add('selected');
            
            const radioInput = this.querySelector('input');
            radioInput.checked = true;
        });
    });
    
    function updateProgress() {
        const progressPercent = ((currentStep - 1) / totalSteps) * 100;
        progressBar.style.width = progressPercent + '%';
        
        steps.forEach(step => {
            const stepNum = parseInt(step.dataset.step);
            step.classList.remove('active', 'completed');
            
            if (stepNum === currentStep) {
                step.classList.add('active');
            } else if (stepNum < currentStep) {
                step.classList.add('completed');
            }
        });
    }

    function goToStep(step) {
        formSteps.forEach(formStep => {
            formStep.classList.remove('active');
        });
    
        document.querySelector(`.form-step[data-step="${step}"]`).classList.add('active');
        
        currentStep = step;
        updateProgress();
    }
    
    nextButtons.forEach(button => {
        button.addEventListener('click', function() {
            const nextStep = parseInt(this.dataset.next);
            
            let isValid = true;
            
            if (nextStep === 2) {
                const fullname = document.getElementById('fullname').value;
                if (!fullname.trim()) {
                    alert('Пожалуйста, введите ваше ФИО');
                    isValid = false;
                }
            } else if (nextStep === 3) {
                const phone = document.getElementById('phone').value;
                if (!phone.trim()) {
                    alert('Пожалуйста, введите ваш номер телефона');
                    isValid = false;
                }
            } else if (nextStep === 4) {
                const positionSelected = document.querySelector('input[name="position"]:checked');
                if (!positionSelected) {
                    alert('Пожалуйста, выберите вашу должность');
                    isValid = false;
                }
            } else if (nextStep === 5) {
                const date = document.getElementById('date').value;
                const timeSelect = document.getElementById('time');
                if (!date) {
                    alert('Пожалуйста, выберите дату');                    
                    isValid = false;
                }

                fetch('/available-times?'+
                    new URLSearchParams({ date: date}).toString()).then(response => {
                        if (!response.ok) throw new Error('Ошибка сети');
                        return response.json();
                    })
                    .then(json => {
                        console.log(json);
                        const availableTimes = json; 
                        timeSelect.innerHTML = '';
                        if (availableTimes.length === 0) {
                            timeSelect.innerHTML = '<option value="" disabled>Нет доступных окон</option>';
                        } else {
                        
                            const defaultOption = document.createElement('option');
                            defaultOption.value = "";
                            defaultOption.disabled = true;
                            defaultOption.selected = true;
                            defaultOption.textContent = "Выберите время";
                            timeSelect.appendChild(defaultOption);
                            
                            
                            availableTimes.forEach(time => {
                                const option = document.createElement('option');
                                option.value = time;
                                option.textContent = time;
                                timeSelect.appendChild(option);
                            });
                            
                            timeSelect.disabled = false;
                        }
                    })
                    .catch(error => {
                        console.error('Ошибка:', error);
                    });
    
            } else if (nextStep === 6) {
                const time = document.getElementById('time').value;
                const empSelect = document.getElementById('employee');
                if (!time) {
                    alert('Пожалуйста, выберите время');
                    isValid = false;
                }
                fetch('/employees_ind').then(response => {
                        if (!response.ok) throw new Error('Ошибка сети');
                        return response.json();
                    })
                    .then(json => {
                        console.log(json);
                        const employees = json; 
                        empSelect.innerHTML = '';
                        if (employees.length === 0) {
                            empSelect.innerHTML = '<option value="" disabled>Нет свободных сотрудников.</option>';
                        } else {
                            const defaultOption = document.createElement('option');
                            defaultOption.value = "";
                            defaultOption.selected = true;
                            defaultOption.textContent = "Не важно";
                            empSelect.appendChild(defaultOption);
                            
                            employees.forEach(emp => {
                                const option = document.createElement('option');
                                option.value = emp;
                                option.textContent = emp;
                                empSelect.appendChild(option);
                            });
                            
                            empSelect.disabled = false;
                        }
                    })
                    .catch(error => {
                        console.error('Ошибка:', error);
                    });
            }
            else if (nextStep === 7) {
                const comment = document.getElementById('comment').value;
            }
            
            if (isValid) {
                goToStep(nextStep);
            }
        });
    });
    
    prevButtons.forEach(button => {
        button.addEventListener('click', function() {
            const prevStep = parseInt(this.dataset.prev);
            goToStep(prevStep);
        });
    });
    
    form.addEventListener('submit', function(e) {
        e.preventDefault();
        
        // Собираем данные формы
        const formData = {
            fio: document.getElementById('fullname').value,
            phone: document.getElementById('phone').value,
            position: document.querySelector('input[name="position"]:checked')?.value || '',
            psyholog: document.getElementById('employee')?.value || '',
            type: "Индивидуальная консультация",
            comment:document.getElementById('comment').value,
            meet_date: document.getElementById('date').value,
            meet_time: document.getElementById('time').value
        };
        

        // Отправляем POST-запрос
        fetch('/event', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData)
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Ошибка сети');
            }
            return response.json();
        })
        .then(data => {
            console.log('Успешно:', data);
            goToStep(8);
        })
        .catch(error => {
            console.error('Ошибка:', error);
            alert('Произошла ошибка при отправке формы');
        });
    });
    
    restartBtn.addEventListener('click', function() {
        
        form.reset();
        
        radioOptions.forEach(opt => opt.classList.remove('selected'));
        
        goToStep(1);
    });
    
    updateProgress();
});