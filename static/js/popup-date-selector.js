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
              confirmButtonColor: "rgb(0,150,0)",
              cancelButtonColor: "rgb(50,50,50)",
              willOpen: () => {
                const data_el = document.getElementById('reservation-date-modal');
                const rp = new DateRangePicker(data_el, {
                    format: "dd-mm-yyyy",
                    minDate: new Date(),
                    showOnFocus: true,
                    clearButton: true,
                    autohide: true,
                    orientation:"top",
                    allowOneSidedRange:"ture",
                  });
              },
              didOpen: () => {
                document.getElementById('start_date').removeAttribute('disabled')
                document.getElementById('end_date').removeAttribute('disabled')
                document.getElementById('reservation-date-modal').style.overflow = "visible"
              },
              preConfirm: () => {
                return [
                  document.getElementById("start_date").value,
                  document.getElementById("end_date").value
                ];
              },
            });
            if (formValues) {
              const start_date = document.getElementById("start_date").value;
              const ende_date = document.getElementById("end_date").value;
              if(start_date === "" || end_date === ""){
                Swal.fire({
                  html: `<b style="color:black;">Invalid Date</b>`,
                  showConfirmButton: false,
                  showCancelButton: true,
                  cancelButtonColor: "rgb(50,50,50)",

                });
              } else if(formValues.dismiss !== Swal.DismissReason.cancel){
                if(formValues.value !== ""){
                  if(c.callback !== undefined){
                    c.callback(formValues);
                  }
                }else{
                  c.callback(false);
                }
              }else{
                c.callback(false);
              }
            }
          }
          return {
            reservation: reservation
          }
        }