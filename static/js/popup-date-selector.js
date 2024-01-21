

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
              // position: "bottom",
              confirmButtonColor: "#28a745",
              cancelButtonColor: "#020202",
              willOpen: () => {
                const data_el = document.getElementById('reservation-date-modal');
                const rp = new DateRangePicker(data_el, {
                    format: "dd-mm-yyyy",
                    showOnFocus: true,
                    clearButton: true,
                    autohide: true,
                    orientation:"top",
                    allowOneSidedRange:"ture",
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