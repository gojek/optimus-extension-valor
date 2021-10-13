local enrich(data) = {
    email: data.email,
    membership: data.membership,
    is_active: data.is_active,
    is_valid: true
};

enrich(user_account)
