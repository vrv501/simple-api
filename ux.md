Personas

Visitor (anonymous): browse/search pets, view pet details, view categories by name.
Authenticated user (buyer/seller): create/manage pets and images they own; add to cart and place orders; track/cancel eligible orders; manage their account.
Admin: manage animal categories (add/update only).
Logistics operator: update order status (requires new admin/operator endpoints).
Animal Categories (admin-only add/update; public read; no delete)

Animal Categories (admin-only add/update; public read; no delete)

Browse category by name (public)
Screen: “Search Categories” (input name) → show result or “not found”.
API: GET /animal-categories (findAnimalCategory)
Handler: apihandler.FindAnimalCategory
Add category (admin)
Screen: “New Category” (name) → create → success toast with created_at.
API: POST /animal-categories (addAnimalCategory)
Handler: apihandler.AddAnimalCategory
Update category (admin)
Screen: “Edit Category” (rename) → save → success toast with updated_at.
API: PUT /animal-categories/{id} (replaceAnimalCategory)
Handler: apihandler.ReplaceAnimalCategory
Delete category (should be disabled)
Current API has DELETE /animal-categories/{id}. Recommend removing this from openapi.yml and deprecating the handler apihandler.DeleteAnimalCategory to align with “cannot be deleted”.

Users

Sign up (public)
Screen: “Create Account” (username, full name, email, phone, password).
API: POST /users (createUser)
Manage account (owner-only)
Screen: “Profile” (edit fields) → PUT /users/{username} (replaceUser)
Delete account (owner-only)
Rule: block if there are undelivered orders.
UX: “Delete Account” → if blocked, show reason and link to “My Orders”.
API: DELETE /users/{username} (deleteUser) with pre-check.

Pets (any authenticated user can list, manage their own pets; images up to 10)

Discover pets (public or auth)
Screen: “Browse Pets” with filters (name, status, tags), pagination via X-Next-Cursor.
API: GET /pets (findPets)
Handler: apihandler.FindPets
Pet detail (public or auth)
Screen: “Pet Detail” with images, price, category, tags; CTA: Add to Cart.
API: GET /pets/{petId} (getPetByID) + GET /pets/{petId}/images/{imageId} (getImageByPetId)
Handlers: apihandler.GetPetByID, apihandler.GetImageByPetId
Create pet (seller)
Screen: “New Pet” (pet JSON + optional photos[≤10]).
API: POST /pets (addPet, multipart)
Handler: apihandler.AddPet
Edit pet (seller)
Screen: “Edit Pet” (fields except status), add/remove images.
API: PUT /pets/{petId} (replacePet, multipart); POST /pets/{petId}/images (uploadPetImage); DELETE /pets/{petId}/images/{imageId} (deletePetImage)
Handlers: apihandler.ReplacePet, apihandler.UploadPetImage, apihandler.DeletePetImage
Enforcement: server ignores/blocks status changes in replacePet.
Delete pet (seller)
Rule: blocked if the pet appears on any placed/processing/shipped/delivered order.
UX: “Delete Pet” → if blocked, show “Pet is part of active or fulfilled orders”.
API: DELETE /pets/{petId} (deletePet)
Handler: apihandler.DeletePets

Cart and Checkout

Cart (client-side or persisted as a draft; no server “cart” resource required)
UX: “Add to Cart” on pet detail → local cart list (pet ids).
Checkout: “Place Order” → confirm address/payment (if applicable) → submit.
API: POST /store/orders (placeOrders) with array of petIds
Handler: apihandler.PlaceOrders

Orders

Order list (user)
Screen: “My Orders” with filters (status, afterDate), pagination via X-Next-Cursor.
API: GET /store/orders (findOrders)
Handler: apihandler.FindOrders
Order detail (user)
Screen: “Order Detail” with per-item shipped/delivered metadata.
API: GET /store/orders/{orderId} (getOrderByID)
Handler: apihandler.GetOrderByID
Cancel/delete order (user)
Rule: allowed only while not shipped (shipped_date is null).
UX: Show “Cancel Order” if eligible; otherwise show disabled state with reason.
API: DELETE /store/orders/{orderId} (deleteOrder)
Handler: apihandler.DeleteOrder
Logistics status updates (operator)
States: placed → processing → shipped → delivered
Actions:
“Mark Processing” when received by partner.
“Mark Shipped” sets shipped_date.
“Mark Delivered” sets delivered date.
API: not present; add admin/operator-only endpoints to update order status (recommend PATCH /store/orders/{orderId}/status).

Suggested status model (minimal churn, clear transitions)

Order.status: placed | processing | shipped | delivered | cancelled
placed: created by user.
processing: received by logistics partner.
shipped: in transit; set shipped_date.
delivered: complete; set delivered_date.
cancelled: user-initiated before shipped (optional; or continue using DELETE to hard-delete).
Keep shipped_date (already present). Add delivered_date. Optionally add cancelled_date.
Enforcement:
Users can DELETE order only while status in {placed, processing} and shipped_date is null (consistent with your rule).
After shipped, no deletion or cancellation.
Pet deletion blocked if referenced by any order with status != cancelled.

Minor schema tweaks (optional, low-friction)

Pet: add breed (string), birthdate (date). Default status remains available, set to sold when delivered. Consider adding reserved to avoid race during checkout, but you can defer that and rely on order placement constraints.
Categories: remove DELETE operation from openapi.yml and deprecate apihandler.DeleteAnimalCategory.
Orders: add delivered_date and cancelled status (and maybe delivered date in response schema). Expose operator-only status update endpoints.
Navigation summary

Public:
Browse Pets → View Pet → Add to Cart → Sign Up/Sign In → Checkout
Search Category by Name
Auth user (seller/buyer):
My Pets: Create/Edit/Delete, Manage Images
My Orders: List/Detail, Cancel/Delete if eligible
Profile: Edit/Delete (blocked with undelivered orders)
Admin:
Categories: Create/Update
Operator:
Orders: Work queue, Mark Processing/Shipped/Delivered
This plan keeps routes stable by:

Using existing operations for most flows in openapi.yml.
Moving “cart” to the client and relying on the existing PlaceOrders array input.
Enforcing business rules in handlers (not new routes), except for operator status updates which need dedicated endpoints.
