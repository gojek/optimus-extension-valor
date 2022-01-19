local evaluate(resource, definition, previous) =
    local membership_dict = definition['memberships'];
    local membership_id = std.toString(resource.membership_id);
    local current_membership = membership_dict[membership_id];
    {
        email: resource.email,
        membership: current_membership.name,
        is_active: resource.is_active
    };
