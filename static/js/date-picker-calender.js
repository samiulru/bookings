const elem = document.getElementById('reservation-date');
        const rangepicker = new DateRangePicker(elem, {
            format: "dd-mm-yyyy",
            showOnFocus: true,
            clearButton: true,
            autohide: true,
            allowOneSidedRange:"ture",
        });