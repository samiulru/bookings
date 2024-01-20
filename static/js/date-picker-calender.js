const elem = document.getElementById('reservation-date');
        const rangepicker = new DateRangePicker(elem, {
            format: "dd-mm-yyyy",
            autohide: true,
        });