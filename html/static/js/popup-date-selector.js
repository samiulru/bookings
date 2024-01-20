
        let attention = Prompt();
        let el = document.getElementById('submitBtn').addEventListener("click", function() {
                    let html = `
                <form id="reservation-date-modal" action="" method="get" novalidate>
                  <h4>Select Dates</h4>            
                  <input type="text" class="custom-popup-form" name="arrival" id="arrival" placeholder="Arrival" autocomplete="off" disabled required>
                  <input type="text" class="custom-popup-form" name="deperture" id="deperture" placeholder="Deperture" autocomplete="off" disabled required>
                </form>`
            attention.reservation({msg: html})
        })

        function Prompt(){
          async function reservation(c){
            const {
              title = "",
              msg = "",
            } = c;
            const { value: formValues } = await Swal.fire({
              title: title,
              html: msg,
              focusConfirm: false,
              showCancelButton: true,
              customClass:"custom-swal-size",
              position: "bottom",
              confirmButtonColor: "#28a745",
              cancelButtonColor: "#020202",
              okButtonColor: "#28a745",
              willOpen: () => {
                const data_el = document.getElementById('reservation-date-modal');
                const rp = new DateRangePicker(data_el, {
                    format: 'dd-mm-yyyy',
                    showOnFocus: true,
                    clearButton: true,
                    autohide: true,
                    orientation:'top',
                  });
              },
              didOpen: () => {
                document.getElementById('arrival').removeAttribute('disabled')
                document.getElementById('deperture').removeAttribute('disabled')
                document.getElementById('reservation-date-modal').style.overflow = "visible"
              },
              preConfirm: () => {
                return [
                  document.getElementById("arrival").value,
                  document.getElementById("deperture").value
                ];
              },
            });
            if (formValues) {
              const arrival_date = document.getElementById("arrival").value;
              const deperture_date = document.getElementById("deperture").value; 
              if(arrival_date == "" || deperture_date == ""){
                Swal.fire({
                  html: `<b style="color:black;">Invalid Date</b>`,
                  showConfirmButton: false,
                  showCancelButton: true,
                  cancelButtonColor: "#020202",
                  
                });
              } else {
                Swal.fire({
                title: "Your Selected Date",
                text: "From " + arrival_date + " To " + deperture_date,
                confirmButtonColor: "#28a745",
              });
              }
            }
          }
          return {
            reservation: reservation
          }
        }