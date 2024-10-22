function editRecord(id) {
    console.log(id)
    // AJAX request to fetch the record data
    $.ajax({
        url: "/records/" + id,
        type: "GET",
        success: function (response) {
            // Handle successful response
            var record = response;
            console.log('success')
            populateFormFields(record);
        },
        error: function (xhr, status, error) {
            // Handle error
            console.error("Error fetching record:", error);
        }
    });
}

function populateFormFields(record) {
    // Populate form fields with record data
    document.getElementById("name").value = record.Name;
    document.getElementById("phone").value = record.Phone;
    document.getElementById("email").value = record.Email;
    document.getElementById("message").value = record.Message;
    document.getElementById("id").value = record.ID; // Set the ID of the record to be updated
    document.getElementById("updateForm").style.display = "block"; // Display the form
}

function deleteRecord(id) {
    // AJAX request to delete the record
    $.ajax({
        url: "/records/" + id,
        type: "DELETE",
        success: function (response) {
            // Handle successful response
            console.log("Record deleted successfully");
        },
        error: function (xhr, status, error) {
            // Handle error
            console.error("Error deleting record:", error);
        }
    });
}
